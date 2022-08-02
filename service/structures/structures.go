package structures

// Shared
type RepositoryParameter struct {
	Description string
	Default     string
}

type CalzoneConfigSystem struct {
	Memory int
	Cpus   int
}

type CalzoneConfigVolume struct {
	Size string
}

type CalzoneConfigDeployment struct {
	Public  bool
	Image   string
	Volumes []string
}

type CalzoneConfig struct {
	System      CalzoneConfigSystem
	Volumes     map[string]CalzoneConfigVolume
	Deployments map[string]CalzoneConfigDeployment
}

// BeginInstall
type BeginInstallArgs struct {
	Calzone string
}

type BeginInstallResponse struct {
	Parameters map[string]RepositoryParameter
}

// CompleteInstall
type CompleteInstallArgs struct {
	Calzone    string
	Parameters map[string]string
}

type CompleteInstallResponse struct {
	Port         int
	MessageQueue []string
}
