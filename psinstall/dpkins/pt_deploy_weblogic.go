package dpkins

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

func WeblogicInstall() error {
	destination := config.Deploy_location
	//destination := getCfgKeyVal("deploy_location")
	archives_dest := config.Archives_dest
	if notExists(filepath.Join(destination, "bea", "weblogic")) {
		err := os.MkdirAll(filepath.Join(destination, "bea"), 0755)
		if err != nil {
			Error.Fatalln(err)
		}
	}
	err = os.MkdirAll(filepath.Join(destination, "bea", "wl"), 0755)
	if err != nil {
		Error.Fatalln(err)
	}
	source := GetArchFile(archives_dest, "weblogic")
	if archiver.TarGz.Match(source) {
		err = archiver.TarGz.Open(source, filepath.Join(destination, "bea", "wl"))
		if err != nil {
			Error.Fatalf("%s => file extract to the folder [ %s ] failed! \n Error message : %v", source, filepath.Join(destination, "bea", "wl"), err)
		}
		Info.Printf("%s ==> extracted successfully to %s\n", source, filepath.Join(destination, "bea", "wl"))
	}
	oracle_home_name := "OraWL1213Home"
	inventory_location := config.Oracle_inventory_location
	wl_ins := filepath.Join(destination, "bea", "wl", "pasteBinary.sh")
	inv_cmd_opt := fmt.Sprintf(`%s/oraInst.loc`, inventory_location)
	clone_cmdf := []string{
		"-javaHome",
		//filepath.Join(getCfgKeyVal("deploy_location"), "java"),
		filepath.Join(destination, "java"),
		"-archiveLoc",
		filepath.Join(destination, "bea", "wl", "pt-weblogic-copy.jar"),
		"-targetMWHomeLoc",
		filepath.Join(destination, "bea", "weblogic"),
		"-targetOracleHomeName",
		oracle_home_name,
		"-invPtrLoc",
		inv_cmd_opt,
		"-executeSysPrereqs",
		"false",
		"-silent",
		"true",
		"-logDirLoc",
		filepath.Join(destination, "bea", "wl"),
	}

	oraInst := []byte(fmt.Sprintf("inventory_loc=/opt/oracle/psft/db/oraInventory\ninst_group=%s", inventory_location))
	err := ioutil.WriteFile(inv_cmd_opt, oraInst, 0755)
	if err != nil {
		Error.Fatalf("%s => unable to create it. Error is %v", inv_cmd_opt, err)
	}

	err = execute(wl_ins, clone_cmdf)
	/*byt, err := exec.Command(wl_ins, clone_cmdf...).CombinedOutput()
	if err != nil {
		Info.Println(string(byt))
		Error.Fatal(err)
	}
	Info.Println(string(byt))*/

	if err != nil {
		return err
	}
	return nil
}
