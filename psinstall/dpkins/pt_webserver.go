package dpkins

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func check_domain_exists(domain_dir, domain_name string) bool {
	Info.Printf("Checking if WebServer domain %s already installed", domain_name)
	domain_exists := true
	domain_path := filepath.Join(domain_dir, "webserv", domain_name)
	//config_file := filepath.Join(domain_path, "config", "config.xml")
	if notExists(filepath.Join(domain_dir, "webserv", domain_name)) {
		domain_exists = false
		if notExists(filepath.Join(domain_path, "config", "config.xml")) {
			domain_exists = false
		}
	}
	return domain_exists
}

func Weblogic_domain_create() error {
	err := create_rep_file("/tmp/resp_fl.txt", &config)
	if err != nil {
		Error.Fatalln(err)
	}
	err = os.Chdir(filepath.Join(config.Deploy_location, "pshome8.55.12", "setup", "PsMpPIAInstall"))
	if err != nil {
		Error.Fatalln(err)
	}

	cmd1 := []string{
		"-i",
		"silent",
		"-DRES_FILE_PATH=/tmp/resp_fl.txt",
		"-tempdir",
		filepath.Join(config.Deploy_location, "tmp"),
	}
	err = execute("./setup.sh", cmd1)
	//byt, err := exec.Command("./setup.sh", cmd1...).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func create_rep_file(fl string, config *tomlConfig) error {
	flH, err := os.Create(fl)
	if err != nil {
		Error.Fatalln(err)
	}
	writer := bufio.NewWriter(flH)
	defer flH.Close() //fmt.Printf("Weblogic Deploy_loc: %s\n", config.Weblogic.Deploy_loc)
	for _, item := range config.Weblogic {
		fmt.Fprintf(writer, "PS_HOME=%s\n", filepath.Join(config.Deploy_location, "pshome8.55.12"))
		fmt.Fprintln(writer, "DOMAIN_NAME=peoplesoft")
		fmt.Fprintln(writer, "DOMAIN_TYPE=NEW_DOMAIN")
		fmt.Fprintln(writer, "INSTALL_ACTION=CREATE_NEW_DOMAIN")
		fmt.Fprintln(writer, "INSTALL_TYPE=SINGLE_SERVER_INSTALLATION")
		fmt.Fprintf(writer, "SERVER_TYPE=%s\n", item.Webserver_type)
		fmt.Fprintf(writer, "BEA_HOME=%s\n", filepath.Join(config.Deploy_location, "bea", "weblogic"))
		fmt.Fprintf(writer, "PS_CFG_HOME=%s\n", config.Ps_config_home)
		fmt.Fprintf(writer, "USER_ID=%s\n", item.Webserver_admin_user)
		fmt.Fprintf(writer, "USER_PWD=%s\n", item.Webserver_admin_user_pwd)
		fmt.Fprintf(writer, "USER_PWD_RETYPE=%s\n", item.Webserver_admin_user_pwd)
		fmt.Fprintf(writer, "WEBSITE_NAME=%s\n", item.Site_name)
		fmt.Fprintf(writer, "HTTP_PORT=%d\n", (item.Webserver_http_port))
		fmt.Fprintf(writer, "HTTPS_PORT=%d\n", (item.Webserver_https_port))
		fmt.Fprintf(writer, "PSSERVER=%s\n", item.Appserver_connStr)
		fmt.Fprintf(writer, "WEB_PROF_NAME=%s\n", item.Profile_name)
		fmt.Fprintf(writer, "WEB_PROF_USERID=%s\n", item.Profile_user)
		fmt.Fprintf(writer, "WEB_PROF_PWD=%s\n", item.Profile_passwd)
		fmt.Fprintf(writer, "WEB_PROF_PWD_RETYPE=%s\n", item.Profile_passwd)
		fmt.Fprintf(writer, "IGW_USERID=%s\n", item.Gateway_user)
		fmt.Fprintf(writer, "IGW_PWD=%s\n", item.Gateway_passwd)
		fmt.Fprintf(writer, "IGW_PWD_RETYPE=%s\n", item.Gateway_passwd)
		fmt.Fprintf(writer, "APPSRVR_CONN_PWD=%s\n", config.Domain_conn_pwd)
		fmt.Fprintf(writer, "APPSRVR_CONN_PWD_RETYPE=%s\n", config.Domain_conn_pwd)
		fmt.Fprintf(writer, "REPORTS_DIR=%s\n", item.Report_repository_dir)
		//response_file.close
		//  File.chmod(0755, response_file_path)
		//  return response_file_path

		//authtoken_domain = resource[:auth_token_domain]
		return writer.Flush()

	}
	return nil
}
