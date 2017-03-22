package main

import (
	"flag"

	"github.com/vkr/psinstall/dpkins"
)

var (
	//archloc = flag.String("archive_dest", "", "Provide the DPK Archive files path.")
	//dplloc  = flag.String("deploy_loc", "", "Path to deploy PeopleTools AppHome & Tuxedo.")
	//domName = flag.String("domain", "", "Enter the Domain / Database Name.")
	appIns = flag.String("installApp", "n", "Only Application Install --> true/false. Defaults to false")
	pshIns = flag.String("installPShome", "n", "Only PShome Install --> true/false. Defaults to false")
	tuxIns = flag.String("installTuxedo", "n", "Only Tuxedo Install --> true/false. Defaults to false")
	webIns = flag.String("installWeblogic", "n", "Only Weblogic Install --> true/false. Defaults to false")
	webDom = flag.String("installWeblogicDomain", "n", "Only Weblogic Domain Install --> true/false. Defaults to false")
	err    error
)

func main() {
	flag.Parse()

	/*if *archloc == "" && *dplloc == "" {
		flag.PrintDefaults()
		log.Println(flag.Args())
		log.Fatalf("Not enough Arcguments : %v", flag.Args())
	}*/

	//	psH, psA, tux := dpkins.GetArchFile(*archloc)

	switch {
	case *appIns == "y" && *pshIns == "n" && *tuxIns == "n":
		dpkins.Info.Println("\n\t ****Working on PSApphome Deploy! ****")
		err := dpkins.ApphomeCreate()
		if err != nil {
			dpkins.Error.Fatalf("PSApphome Deploy Failed => %v", err)
		} else {
			dpkins.Info.Println("\n\t ****PSApphome Deploy Completed Successfully ****")
		}
	case *pshIns == "y" && *appIns == "n" && *tuxIns == "n":
		dpkins.Info.Println("\n\t **** Working on PShome Deploy! ****")
		err := dpkins.PshomeCreate()
		if err != nil {
			dpkins.Error.Fatalf("PShome Deploy Failed => %v", err)
		} else {
			dpkins.Info.Println("\n\t ****PShome Deploy Completed Successfully ****")
		}
	case *pshIns == "n" && *appIns == "n" && *tuxIns == "y":
		dpkins.Info.Println("\n\t **** Working on Tuxedo Deploy! ****")
		err := dpkins.TuxedoCreate()
		if err != nil {
			dpkins.Error.Fatalf("TuxHome Creation Failed => %v", err)
		} else {
			dpkins.Info.Println("\n\t ****TuxHome Deploy Completed Successfully ****")
		}
	case *pshIns == "n" && *appIns == "n" && *tuxIns == "n" && *webIns == "y":
		dpkins.Info.Println("\n\t **** Working on Weblogic Deploy! ****")
		err := dpkins.WeblogicInstall()
		if err != nil {
			dpkins.Error.Fatalf("Weblogic Deploy Failed => %v", err)
		} else {
			dpkins.Info.Println("\n\t ****Weblogic Deploy Completed Successfully ****")
		}
	case *pshIns == "n" && *appIns == "n" && *tuxIns == "n" && *webIns == "n" && *webDom == "y":
		dpkins.Info.Println("\n\t **** Working on Weblogic Domain Creation ! ****")
		err := dpkins.Weblogic_domain_create()
		if err != nil {
			dpkins.Error.Fatalf("\n\t **** Weblogic Domain Creation Failed => %v", err)
		} else {
			dpkins.Info.Println("\n\t **** Weblogic Domain Creation Completed Successfully ****")
		}
	default:
		dpkins.Info.Println("No matching install criteria provided! => ")
		dpkins.Info.Println(dpkins.PrintUsage())
	}
}
