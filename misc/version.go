package misc

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	gitRef, gitSHA string
)

func SetGitRevisionDetails(r, s string) {
	gitRef = r
	gitSHA = s
}

func DetectVersionDuringRuntime() {
	logrus.Info("starting version details during runtime!")
	logrus.Debug("trying to get details from build debug details")

	bi, ok := debug.ReadBuildInfo()
	if ok {
		logrus.Debug("build info was detected, so using it...")
		gitRef = bi.Main.Version
		for _, s := range bi.Settings {
			if s.Key == "vcs.revision" {
				gitSHA = s.Value
				break
			}
		}
		if gitRef != "" && gitSHA != "" {
			logrus.Debugf("detected git ref [%s] and git sha [%s] from build info", gitRef, gitSHA)
			return
		}
	}

	_, b, _, _ := runtime.Caller(0)
	projectPath := path.Join(filepath.Dir(b), "..")
	logrus.Debugf("could not get version details from build, will try to get it from local git at [%s]", projectPath)

	cmd := exec.Command("git", "show-ref", "--hash=8", "HEAD")
	cmd.Dir = projectPath
	stdout, err := cmd.Output()
	if err != nil {
		logrus.Debugf("could not detect local git revision: %s", err)
		return
	}
	gitSHA = strings.TrimSpace(string(stdout))

	cmd = exec.Command("git", "symbolic-ref", "HEAD")
	cmd.Dir = projectPath
	stdout, err = cmd.Output()
	if err != nil {
		logrus.Debugf("could not detect local git branch: [%s]", err)
		return
	}
	gitRef = fmt.Sprintf("%s:%s", os.Getenv("USER"), strings.TrimSpace((string(stdout))))

	logrus.Infof("detected details from local GIT: ref [%s] sha [%s]", gitRef, gitSHA)
}

func GetGitRef() string {
	return gitRef
}

func GetGitSHA() string {
	return gitSHA
}

func GetVersion() string {
	// when face with refs/tags
	if strings.HasPrefix(gitRef, "refs/tags") {
		// local branch
		return strings.Split(gitRef, "/")[2] + ":" + gitSHA
	} else if gitRef != "" {
		return gitRef + ":" + gitSHA
	}
	return "VersionNotDetected"
}
