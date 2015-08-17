package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	path    string
	out     string
	device  string
	browser string
	list    bool
)

const (
	version     = "0.1.0"
	aresVersion = "1.1.0-j393-k"
)

func init() {
	//chrome != google-chrome
	flag.StringVar(&path, "path", ".", "declare a path to app folder")
	flag.StringVar(&path, "p", ".", "declare a path to app folder(shorthand)")
	flag.StringVar(&out, "output", "", "declare a output path, default put ipk in build/ folder")
	flag.StringVar(&out, "o", "", "declare a output path, default put ipk in build/ folder(shorthand)")
	flag.StringVar(&device, "device", "webOs", "declare a device to deploy")
	flag.StringVar(&device, "d", "webOs", "declare a device to deploy(shorthand)")
	flag.StringVar(&browser, "browser", "", "declare a default browser for run inspector")
	flag.StringVar(&browser, "b", "", "declare a default browser for run inspector(shorthand)")
	flag.BoolVar(&list, "list", false, "set true to get device list in json format")
	flag.BoolVar(&list, "l", false, "set true to get device list in json format(shorthand)")
}

func main() {
	flag.Parse()

	curVersion := Version()

	curVersion = strings.Split(strings.Split(curVersion, " ")[1], "\n")[0]
	if curVersion != aresVersion {
		fmt.Println("Warning! Ares version is different. ares-deploy using '1.1.0-j393-k' of ares-setup-device.\n Some functional may not work.")
	}
	if list == true {
		fmt.Println(ListDevice())
	} else {
		appId, filename := ParseInfo(path)
		if out == "" {
			out = fmt.Sprint(path, "/build")
		}
		err, host := Location(device)
		if err != nil {
			log.Fatalln(err)
		}

		Package(path, out)
		Install(device, out, filename)
		Launch(device, appId)

		LaunchBrowser(browser, host)
		fmt.Println("Finish")
	}
}

func ParseInfo(p string) (id, filename string) {
	var jsonStruct interface{}
	filePath := fmt.Sprint(p, "/appinfo.json")
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("appinfo.json not found.")
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &jsonStruct)
	if err != nil {
		log.Fatal(err)
	}
	id, filename = PackageName(jsonStruct)
	return
}

func PackageName(j interface{}) (id, packageName string) {
	g := j.(map[string]interface{})
	if g["id"] == nil {
		log.Fatal("Please specify id field in 'package.json' file.")
	}
	id = g["id"].(string)
	version := g["version"]
	if version != nil {
		packageName = fmt.Sprintf("%s_%s_all.ipk", id, version.(string))
	} else {
		packageName = fmt.Sprintf("%s.ipk", id)
	}
	return
}

func Package(p, o string) {
	cmd := exec.Command("ares-package", p, "-o", o)
	Output(cmd)
}

func Install(device, out, filename string) {
	str := fmt.Sprint(out, "/", filename)
	cmd := exec.Command("ares-install", "--device", device, str)
	Output(cmd)
}

func Launch(device, id string) {
	cmd := exec.Command("ares-launch", "--device", device, id)
	Output(cmd)
}

func Output(cmd *exec.Cmd) {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func Version() (version string) {
	out, err := exec.Command("ares-setup-device", "-V").Output()
	if err != nil {
		log.Fatalln(err)
	}
	version = string(out)
	return
}

func ListDevice() (list string) {
	out, err := exec.Command("ares-setup-device", "-F").Output()
	if err != nil {
		log.Fatalln(err)
	}
	list = string(out)
	return
}

func LaunchBrowser(browser, host string) {
	location := fmt.Sprint("http://", host, ":9998")
	var cmd *exec.Cmd
	//TODO: cut into a function
	if browser == "" {
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/C", "start", location)
		case "linux":
			cmd = exec.Command("xdg-open", location)
		case "darwin":
			cmd = exec.Command("open", location)
		default:
			fmt.Println("Unsupported platform")
			return
		}
	} else {
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/C", "start", browser, location)
		case "linux":
			cmd = exec.Command(browser, location)
		case "darwin":
			cmd = exec.Command("open", "-a", browser, location)
		default:
			fmt.Println("Unsupported platform")
			return
		}
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Please check that your browser specified in a PATH")
		log.Fatalln(err)
	}
}

func Location(device string) (err error, address string) {
	var list []map[string]interface{}
	cmd := exec.Command("ares-setup-device", "-F")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(out, &list)
	if err != nil {
		log.Fatal(err)
	}
	for _, k := range list {
		if k["name"] == device {
			info := k["deviceinfo"].(map[string]interface{})
			address = info["ip"].(string)
			break
		}
	}
	if address == "" {
		err = fmt.Errorf("Can't find device with such name: %s", device)
	}
	return err, address
}
