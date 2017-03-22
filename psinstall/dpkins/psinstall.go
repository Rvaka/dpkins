package dpkins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/vkr/psinstall/archiver"
)

func init() {
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
}

func notExists(fd string) bool {
	_, err = os.Stat(fd)
	if os.IsNotExist(err) {
		return true
	} else {
		return false
	}
	return true

}

/*func getCfgKeyVal(keyName string) string {
	value := fmt.Sprint(viper.Get(keyName))
	if value == "nil" {
		Info.Printf("%s => key cannot be nil.", keyName)
		os.Exit(1)
	}
	return value
}*/

func ApphomeCreate() error {
	destination := config.Deploy_location
	archives_dest := config.Archives_dest
	domName := config.App_domain
	if notExists(destination) {
		Error.Fatalf("%s => Deploy Location doesn't exist : ", destination)
	}
	if notExists(filepath.Join(destination, "hcm92", domName)) {
		os.Mkdir(filepath.Join(destination, "hcm92", domName), 0777)
	}
	source := GetArchFile(archives_dest, "psapphome")
	if archiver.TarGz.Match(source) {
		err = archiver.TarGz.Open(source, filepath.Join(destination, "hcm92", domName))
		if err != nil {
			Error.Fatalf("%s => file extract to the folder [ %s ] failed! ", source, filepath.Join(destination, "hcm92", domName))
		}
	}

	return nil
}

func TuxedoCreate() error {
	destination := config.Deploy_location
	archives_dest := config.Archives_dest
	if notExists(filepath.Join(destination, "bea", "tuxedo")) {
		err := os.MkdirAll(filepath.Join(destination, "bea"), 0755)
		if err != nil {
			Error.Fatalln(err)
		}
	}
	source := GetArchFile(archives_dest, "tuxedo")
	if archiver.TarGz.Match(source) {
		err = archiver.TarGz.Open(source, filepath.Join(destination, "bea", "tuxedo"))
		if err != nil {
			Error.Fatalf("%s => file extract to the folder [ %s ] failed! \n Error message : %v", source, filepath.Join(destination, "bea", "tuxedo"), err)
		}
		Info.Printf("%s ==> extracted successfully to %s\n", source, filepath.Join(destination, "bea", "tuxedo"))
	}
	tuxedo_home_name := "OraTux1213Home"
	inventory_location := config.Oracle_inventory_location
	tuxedo_home := filepath.Join(config.Deploy_location, "bea", "tuxedo")
	clone_cmd := "runInstaller"
	//cmd_suffix := "\"
	inv_cmd_opt := fmt.Sprintf(`%s/oraInst.loc`, inventory_location)
	clone_cmd_path := filepath.Join(tuxedo_home, "oui", "bin", clone_cmd)
	clone_cmdf := []string{
		"-silent",
		"-clone",
		"-waitforcompletion",
		"-nowait",
		"-invPtrLoc",
		inv_cmd_opt,
		fmt.Sprintf("ORACLE_HOME=%s", tuxedo_home),
		fmt.Sprintf("ORACLE_HOME_NAME=%s", tuxedo_home_name),
		"TLISTEN_PASSWORD=password",
	}

	oraInst := []byte(fmt.Sprintf("inventory_loc=/opt/oracle/psft/db/oraInventory\ninst_group=%s", inventory_location))
	err := ioutil.WriteFile(inv_cmd_opt, oraInst, 0755)
	if err != nil {
		Error.Fatalf("%s => unable to create it. Error is %v", inv_cmd_opt, err)
	}
	err = execute(clone_cmd_path, clone_cmdf)
	if err != nil {
		return err
	}
	return nil
}

func PshomeCreate() error {
	destination := config.Deploy_location
	archives_dest := config.Archives_dest
	if notExists(destination) {
		Error.Fatalf("%s => Deploy Location doesn't exist : ", destination)
	}
	source := GetArchFile(archives_dest, "pshome")
	var src = true
	var targ string
	switch src {
	case strings.Contains(strings.ToLower(source), "pshome"):
		_, fl := filepath.Split(source)
		targ = strings.Split(fl, "-")[1]
		targ = strings.TrimSuffix(targ, ".tgz")
		targ = filepath.Join(destination, targ)
	default:
		Info.Println("Inside Case default")
	}

	if archiver.TarGz.Match(source) {
		err = archiver.TarGz.Open(source, targ)
		if err != nil {
			Error.Fatalf("[ %s ] => file extract to the folder [ %s ] failed! \n Error Message : %v", source, targ, err)
		}
	}

	if err == os.Rename(filepath.Join(targ, "ORA", "sqr"), filepath.Join(targ, "sqr")) && err != nil {
		Error.Fatalf("mv %s ==> %s failed! \n Error Message : %v", filepath.Join(targ, "ORA", "sqr"), filepath.Join(targ, "sqr"), err)
	}
	if err == os.Rename(filepath.Join(targ, "ORA", "bin", "sqr"), filepath.Join(targ, "bin", "sqr")) && err != nil {
		Error.Fatalf("mv %s ==> %s failed! \n Error Message : %v", filepath.Join(targ, "ORA", "sqr", "bin"), filepath.Join(targ, "bin", "sqr"), err)
	}
	if err == os.Rename(filepath.Join(targ, "ORA", "scripts"), filepath.Join(targ, "scripts")) && err != nil {
		Error.Println("Error message : ", err)
		Error.Fatalf("mv %s ==> %s failed! ", filepath.Join(targ, "ORA", "scripts"), filepath.Join(targ, "scripts"))
	}
	if err == os.Rename(filepath.Join(targ, "ORA", "setup", "psdb.sh"),
		filepath.Join(targ, "setup", "psdb.sh")) && err != nil {
		Error.Println("Error message : ", err)
		Error.Fatalf("mv %s ==> %s failed! ", filepath.Join(targ, "ORA", "setup", "psdb.sh"), filepath.Join(targ, "setup", "psdb.sh"))
	}
	if err == os.Rename(filepath.Join(targ, "ORA", "psconfig.sh"), filepath.Join(targ, "psconfig.sh")) && err != nil {
		Error.Println("Error message : ", err)
		Error.Fatalf("mv %s ==> %s failed! ", filepath.Join(targ, "ORA", "psconfig.sh"), filepath.Join(targ, "psconfig.sh"))
	}
	if err == os.Rename(filepath.Join(targ, "ORA", "setup", "dbcodes.pt"),
		filepath.Join(targ, "setup", "dbcodes.pt")) && err != nil {
		Error.Println("Error message : ", err)
		Error.Fatalf("mv %s ==> %s failed! ", filepath.Join(targ, "ORA", "setup", "dbcodes.pt"), filepath.Join(targ, "setup", "dbcodes.pt"))
	}

	if err == os.Rename(filepath.Join(targ, "ORA", "peopletools.properties"),
		filepath.Join(targ, "peopletools.properties")) && err != nil {
		Error.Println("Error message : ", err)
		Error.Fatalf("mv %s ==> %s failed! ", filepath.Join(targ, "ORA", "peopletools.properties"), filepath.Join(targ, "peopletools.properties"))
	}

	return nil
}

func GetArchFile(archives, srchstr string) (archive_file string) {
	if notExists(archives) {
		Error.Fatalf("%s => Archives Folder path is not valid : ", archives)
	}

	filepath.Walk(archives, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			Error.Print(err)
			return nil
		}
		if strings.HasSuffix(path, ".tgz") {
			if strings.Contains(strings.ToLower(path), srchstr) {
				archive_file = path
			}
		}
		return nil
	})
	if archive_file == "nil" {
		Error.Fatalf("%s archive file doesn't exist in the archive folder %s", srchstr, archives)
	}
	return archive_file
}
