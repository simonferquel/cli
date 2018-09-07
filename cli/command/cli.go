package command

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/config"
	cliconfig "github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/connhelper"
	cliflags "github.com/docker/cli/cli/flags"
	manifeststore "github.com/docker/cli/cli/manifest/store"
	registryclient "github.com/docker/cli/cli/registry/client"
	"github.com/docker/cli/cli/trust"
	dopts "github.com/docker/cli/opts"
	clitypes "github.com/docker/cli/types"
	"github.com/docker/docker/api"
	"github.com/docker/docker/api/types"
	registrytypes "github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	dcontext "github.com/docker/docker/client/context"
	"github.com/docker/docker/pkg/contextstore"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/theupdateframework/notary"
	notaryclient "github.com/theupdateframework/notary/client"
	"github.com/theupdateframework/notary/passphrase"
)

const (
	// ContextDockerHost is the reported context when DOCKER_HOST is used
	ContextDockerHost = "<DOCKER_HOST>"
)

// Streams is an interface which exposes the standard input and output streams
type Streams interface {
	In() *InStream
	Out() *OutStream
	Err() io.Writer
}

// Cli represents the docker command line client.
type Cli interface {
	Client() client.APIClient
	Out() *OutStream
	Err() io.Writer
	In() *InStream
	SetIn(in *InStream)
	ConfigFile() *configfile.ConfigFile
	ServerInfo() ServerInfo
	ClientInfo() ClientInfo
	NotaryClient(imgRefAndAuth trust.ImageRefAndAuth, actions []string) (notaryclient.Repository, error)
	DefaultVersion() string
	ManifestStore() manifeststore.Store
	RegistryClient(bool) registryclient.RegistryClient
	ContentTrustEnabled() bool
	NewContainerizedEngineClient(sockPath string) (clitypes.ContainerizedClient, error)
	ContextStore() contextstore.Store
	CurrentContext() string
}

// DockerCli is an instance the docker command line client.
// Instances of the client can be returned from NewDockerCli.
type DockerCli struct {
	configFile            *configfile.ConfigFile
	in                    *InStream
	out                   *OutStream
	err                   io.Writer
	client                client.APIClient
	serverInfo            ServerInfo
	clientInfo            ClientInfo
	contentTrust          bool
	newContainerizeClient func(string) (clitypes.ContainerizedClient, error)
	contextStore          contextstore.Store
	currentContext        string
}

// DefaultVersion returns api.defaultVersion or DOCKER_API_VERSION if specified.
func (cli *DockerCli) DefaultVersion() string {
	return cli.clientInfo.DefaultVersion
}

// Client returns the APIClient
func (cli *DockerCli) Client() client.APIClient {
	return cli.client
}

// Out returns the writer used for stdout
func (cli *DockerCli) Out() *OutStream {
	return cli.out
}

// Err returns the writer used for stderr
func (cli *DockerCli) Err() io.Writer {
	return cli.err
}

// SetIn sets the reader used for stdin
func (cli *DockerCli) SetIn(in *InStream) {
	cli.in = in
}

// In returns the reader used for stdin
func (cli *DockerCli) In() *InStream {
	return cli.in
}

// ShowHelp shows the command help.
func ShowHelp(err io.Writer) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SetOutput(err)
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}

// ConfigFile returns the ConfigFile
func (cli *DockerCli) ConfigFile() *configfile.ConfigFile {
	return cli.configFile
}

// ServerInfo returns the server version details for the host this client is
// connected to
func (cli *DockerCli) ServerInfo() ServerInfo {
	return cli.serverInfo
}

// ClientInfo returns the client details for the cli
func (cli *DockerCli) ClientInfo() ClientInfo {
	return cli.clientInfo
}

// ContentTrustEnabled returns whether content trust has been enabled by an
// environment variable.
func (cli *DockerCli) ContentTrustEnabled() bool {
	return cli.contentTrust
}

// BuildKitEnabled returns whether buildkit is enabled either through a daemon setting
// or otherwise the client-side DOCKER_BUILDKIT environment variable
func BuildKitEnabled(si ServerInfo) (bool, error) {
	buildkitEnabled := si.BuildkitVersion == types.BuilderBuildKit
	if buildkitEnv := os.Getenv("DOCKER_BUILDKIT"); buildkitEnv != "" {
		var err error
		buildkitEnabled, err = strconv.ParseBool(buildkitEnv)
		if err != nil {
			return false, errors.Wrap(err, "DOCKER_BUILDKIT environment variable expects boolean value")
		}
	}
	return buildkitEnabled, nil
}

// ManifestStore returns a store for local manifests
func (cli *DockerCli) ManifestStore() manifeststore.Store {
	// TODO: support override default location from config file
	return manifeststore.NewStore(filepath.Join(config.Dir(), "manifests"))
}

// RegistryClient returns a client for communicating with a Docker distribution
// registry
func (cli *DockerCli) RegistryClient(allowInsecure bool) registryclient.RegistryClient {
	resolver := func(ctx context.Context, index *registrytypes.IndexInfo) types.AuthConfig {
		return ResolveAuthConfig(ctx, cli, index)
	}
	return registryclient.NewRegistryClient(resolver, UserAgent(), allowInsecure)
}

// ContextStore returns the context store
func (cli *DockerCli) ContextStore() contextstore.Store {
	return cli.contextStore
}

// CurrentContext returns the current context name
func (cli *DockerCli) CurrentContext() string {
	return cli.currentContext
}

// Initialize the dockerCli runs initialization that must happen after command
// line flags are parsed.
func (cli *DockerCli) Initialize(opts *cliflags.ClientOptions) error {

	cli.configFile = cliconfig.LoadDefaultConfigFile(cli.err)

	var err error
	cli.contextStore, err = contextstore.NewStore(cliconfig.ContextStoreDir())
	if err != nil {
		return err
	}
	cli.client, cli.currentContext, err = NewAPIClientFromFlags(opts.Common, cli.configFile)
	if tlsconfig.IsErrEncryptedKey(err) {
		passRetriever := passphrase.PromptRetrieverWithInOut(cli.In(), cli.Out(), nil)
		newClient := func(password string) (client.APIClient, string, error) {
			opts.Common.TLSOptions.Passphrase = password
			return NewAPIClientFromFlags(opts.Common, cli.configFile)
		}
		cli.client, cli.currentContext, err = getClientWithPassword(passRetriever, newClient)
	}
	if err != nil {
		return err
	}
	var experimentalValue string
	// Environment variable always overrides configuration
	if experimentalValue = os.Getenv("DOCKER_CLI_EXPERIMENTAL"); experimentalValue == "" {
		experimentalValue = cli.configFile.Experimental
	}
	hasExperimental, err := isEnabled(experimentalValue)
	if err != nil {
		return errors.Wrap(err, "Experimental field")
	}
	cli.clientInfo = ClientInfo{
		DefaultVersion:  cli.client.ClientVersion(),
		HasExperimental: hasExperimental,
	}
	cli.initializeFromClient()
	return nil
}

func isEnabled(value string) (bool, error) {
	switch value {
	case "enabled":
		return true, nil
	case "", "disabled":
		return false, nil
	default:
		return false, errors.Errorf("%q is not valid, should be either enabled or disabled", value)
	}
}

func (cli *DockerCli) initializeFromClient() {
	ping, err := cli.client.Ping(context.Background())
	if err != nil {
		// Default to true if we fail to connect to daemon
		cli.serverInfo = ServerInfo{HasExperimental: true}

		if ping.APIVersion != "" {
			cli.client.NegotiateAPIVersionPing(ping)
		}
		return
	}

	cli.serverInfo = ServerInfo{
		HasExperimental: ping.Experimental,
		OSType:          ping.OSType,
		BuildkitVersion: ping.BuilderVersion,
	}
	cli.client.NegotiateAPIVersionPing(ping)
}

func getClientWithPassword(passRetriever notary.PassRetriever, newClient func(password string) (client.APIClient, string, error)) (client.APIClient, string, error) {
	for attempts := 0; ; attempts++ {
		passwd, giveup, err := passRetriever("private", "encrypted TLS private", false, attempts)
		if giveup || err != nil {
			return nil, "", errors.Wrap(err, "private key is encrypted, but could not get passphrase")
		}

		apiclient, ctxName, err := newClient(passwd)
		if !tlsconfig.IsErrEncryptedKey(err) {
			return apiclient, ctxName, err
		}
	}
}

// NotaryClient provides a Notary Repository to interact with signed metadata for an image
func (cli *DockerCli) NotaryClient(imgRefAndAuth trust.ImageRefAndAuth, actions []string) (notaryclient.Repository, error) {
	return trust.GetNotaryRepository(cli.In(), cli.Out(), UserAgent(), imgRefAndAuth.RepoInfo(), imgRefAndAuth.AuthConfig(), actions...)
}

// NewContainerizedEngineClient returns a containerized engine client
func (cli *DockerCli) NewContainerizedEngineClient(sockPath string) (clitypes.ContainerizedClient, error) {
	return cli.newContainerizeClient(sockPath)
}

// ServerInfo stores details about the supported features and platform of the
// server
type ServerInfo struct {
	HasExperimental bool
	OSType          string
	BuildkitVersion types.BuilderVersion
}

// ClientInfo stores details about the supported features of the client
type ClientInfo struct {
	HasExperimental bool
	DefaultVersion  string
}

// NewDockerCli returns a DockerCli instance with IO output and error streams set by in, out and err.
func NewDockerCli(in io.ReadCloser, out, err io.Writer, isTrusted bool, containerizedFn func(string) (clitypes.ContainerizedClient, error)) *DockerCli {
	return &DockerCli{in: NewInStream(in), out: NewOutStream(out), err: err, contentTrust: isTrusted, newContainerizeClient: containerizedFn}
}

// NewAPIClientFromFlags creates a new APIClient from command line flags
func NewAPIClientFromFlags(opts *cliflags.CommonOptions, configFile *configfile.ConfigFile) (client.APIClient, string, error) {
	s, err := contextstore.NewStore(cliconfig.ContextStoreDir())
	if err != nil {
		return nil, "", err
	}
	contextName := resolveContextName(opts, s)
	if contextName == ContextDockerHost {
		cli, err := newAPIClientFromFlagsNoContext(opts, configFile)
		return cli, ContextDockerHost, err
	}
	customHeaders := configFile.HTTPHeaders
	if customHeaders == nil {
		customHeaders = map[string]string{}
	}
	customHeaders["User-Agent"] = UserAgent()
	cli, err := client.NewClientWithOpts(
		client.WithContextStoreOrEnv(cliconfig.ContextStoreDir(), contextName),
		client.WithHTTPHeaders(customHeaders),
	)
	if err != nil {
		return nil, "", err
	}
	return cli, contextName, err
}

func newAPIClientFromFlagsNoContext(opts *cliflags.CommonOptions, configFile *configfile.ConfigFile) (client.APIClient, error) {
	host, err := getServerHost(opts.Hosts, opts.TLSOptions)
	if err != nil {
		return &client.Client{}, err
	}
	var clientOpts []func(*client.Client) error
	helper, err := connhelper.GetConnectionHelper(host)
	if err != nil {
		return &client.Client{}, err
	}
	if helper == nil {
		clientOpts = append(clientOpts, withHTTPClient(opts.TLSOptions))
		clientOpts = append(clientOpts, client.WithHost(host))
	} else {
		clientOpts = append(clientOpts, func(c *client.Client) error {
			httpClient := &http.Client{
				// No tls
				// No proxy
				Transport: &http.Transport{
					DialContext: helper.Dialer,
				},
			}
			return client.WithHTTPClient(httpClient)(c)
		})
		clientOpts = append(clientOpts, client.WithHost(helper.Host))
		clientOpts = append(clientOpts, client.WithDialContext(helper.Dialer))
	}

	customHeaders := configFile.HTTPHeaders
	if customHeaders == nil {
		customHeaders = map[string]string{}
	}
	customHeaders["User-Agent"] = UserAgent()
	clientOpts = append(clientOpts, client.WithHTTPHeaders(customHeaders))

	verStr := api.DefaultVersion
	if tmpStr := os.Getenv("DOCKER_API_VERSION"); tmpStr != "" {
		verStr = tmpStr
	}
	clientOpts = append(clientOpts, client.WithVersion(verStr))

	return client.NewClientWithOpts(clientOpts...)
}

func getServerHost(hosts []string, tlsOptions *tlsconfig.Options) (string, error) {
	var host string
	switch len(hosts) {
	case 0:
		host = os.Getenv("DOCKER_HOST")
	case 1:
		host = hosts[0]
	default:
		return "", errors.New("Please specify only one -H")
	}

	return dopts.ParseHost(tlsOptions != nil, host)
}

func withHTTPClient(tlsOpts *tlsconfig.Options) func(*client.Client) error {
	return func(c *client.Client) error {
		if tlsOpts == nil {
			// Use the default HTTPClient
			return nil
		}

		opts := *tlsOpts
		opts.ExclusiveRootPools = true
		tlsConfig, err := tlsconfig.Client(opts)
		if err != nil {
			return err
		}

		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
				DialContext: (&net.Dialer{
					KeepAlive: 30 * time.Second,
					Timeout:   30 * time.Second,
				}).DialContext,
			},
			CheckRedirect: client.CheckRedirect,
		}
		return client.WithHTTPClient(httpClient)(c)
	}
}

// UserAgent returns the user agent string used for making API requests
func UserAgent() string {
	return "Docker-Client/" + cli.Version + " (" + runtime.GOOS + ")"
}

func resolveContextName(opts *cliflags.CommonOptions, store contextstore.Store) string {
	if opts.Context != "" {
		return opts.Context
	}
	if len(opts.Hosts) > 0 {
		return ContextDockerHost
	}
	if _, present := os.LookupEnv("DOCKER_HOST"); present {
		return ContextDockerHost
	}
	if ctxName, ok := os.LookupEnv(dcontext.DockerContextEnvVar); ok {
		return ctxName
	}
	ctxName := store.GetCurrentContext()
	if ctxName == "" {
		return ContextDockerHost
	}
	return ctxName
}
