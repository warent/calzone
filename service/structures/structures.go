package structures

// Shared
type RepositoryParameter struct {
	Description string
	Default     string
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
	Port int
}
