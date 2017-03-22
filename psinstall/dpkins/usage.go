package dpkins

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Oracle_inventory_location string
	Oracle_inventory_user     string
	Oracle_inventory_group    string
	Deploy_location           string
	Archives_dest             string
	App_domain                string
	Ps_config_home            string
	Domain_conn_pwd           string
	Tuxedo                    TuxedoInfo
	Pshome                    PshomeInfo
	Weblogic                  map[string]Property
}

type TuxedoInfo struct{}

type PshomeInfo struct {
	Ps_home_dir string
}

type Property struct {
	Webserver_type           string
	Webserver_admin_user     string
	Webserver_admin_user_pwd string
	Site_name                string
	Appserver_connStr        string
	Profile_name             string
	Profile_user             string
	Profile_passwd           string
	Gateway_user             string
	Gateway_passwd           string
	Report_repository_dir    string
	Webserver_admin_port     int
	Webserver_http_port      int
	Webserver_https_port     int
}

type AppDomain struct {
	dbname        string
	oprid         string
	oprid_pswd    string
	domain_name   string
	connect_id    string
	connect_pwd   string
	server_name   string
	template_type string
}

var config tomlConfig

var usage = `Usage: psInstall {installPShome|installApp|installTuxedo|installWeblogic|installWeblogicDomain}

  -installPShome
    PSHome archive file will be feeded from the DPK's and
		will be installed in the mentioned location.

  -installApp
	Application archive file will be feeded from DPK's and
	will be installed in the mentioned location.

	-installTuxedo
	Tuxedo archive file will be feeded from DPK's and
	will be installed in the mentioned location.

	-installWeblogic
	Weblogc archive file will be feeded from DPK's and
	will be installed in the mentioned location.

	-installWeblogicDomain
	Weblogc archive file will be feeded from DPK's and
	will be installed in the mentioned location.
	`

func init() {
	if _, err := toml.DecodeFile(`./app.toml`, &config); err != nil {
		fmt.Println(err)
		return
	}

}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func PrintUsage() string {
	return usage
}

func execute(cmd_path string, cmdopts []string) error {
	cmd := exec.Command(cmd_path, cmdopts...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Error.Fatalf("Error is %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		Error.Fatalf("Error is %v", err)
	}
	// Start command
	err = cmd.Start()
	if err != nil {
		Error.Fatalf("Error is %v", err)
	}

	defer cmd.Wait() // Doesn't block

	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	return nil
}
