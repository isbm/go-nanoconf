/*
Nanoconf finder helper.
A convenience class that allows to organise configuration files for the application.
Configuration files might be found in /etc/somewhere.conf, $HOME/.somewhere,
$HOME/.config/yourapp/somewhere.conf, /opt/yourapp/etc/somewhere.conf etc etc etc.

To organise all that "mess", NanoconfFinder is to serve.
*/

package nanoconf

import (
	"os"
	"os/user"
	"path"
)

type NanoconfFinder struct {
	appname     string
	defaultconf string
	names       []string
	paths       []string
}

// Constructor
func NewNanoconfFinder(appname string) *NanoconfFinder {
	nf := new(NanoconfFinder)
	nf.appname = appname
	nf.names = make([]string, 0)
	nf.paths = make([]string, 0)
	return nf
}

// Setup typical locations for the user: /etc, $HOME/.config..., $HOME/.<cfg> etc.
// If user is nil, current user is used.
func (nf *NanoconfFinder) DefaultSetup(usr *user.User) *NanoconfFinder {
	var err error
	if usr == nil {
		usr, err = user.Current()
		if err != nil {
			panic(err)
		}
	}
	nf.paths = append(nf.paths, []string{"/etc", path.Join("/etc", "appname"), usr.HomeDir, path.Join(usr.HomeDir, ".config", nf.appname), "./"}...)
	nf.names = append(nf.names, []string{nf.appname + ".conf", "." + nf.appname}...)
	return nf
}

// FindFirst finds nearest  configuration that matches the search
func (nf *NanoconfFinder) FindFirst() string {
	var pth string
	config := nf.FindAll()
	if len(config) > 0 {
		pth = config[len(config)-1]
	}
	return pth
}

// FindAll finds all possible configuration files that would match the search.
func (nf *NanoconfFinder) FindAll() []string {
	configs := make([]string, 0)
	for _, pth := range nf.paths {
		for _, nm := range nf.names {
			cfgpth := path.Join(pth, nm)
			if fh, err := os.Stat(cfgpth); os.IsExist(err) && !fh.IsDir() {
				configs = append(configs, cfgpth)
			}
		}
	}
	return configs
}

// FindDefault returns assumed confguration that could be default, unless it is pre-set.
func (nf *NanoconfFinder) FindDefault() string {
	if nf.defaultconf != "" {
		return nf.defaultconf
	}
	return path.Join("/etc", nf.appname+".conf")
}

func (nf *NanoconfFinder) SetDefaultConfig(cfgpath string) *NanoconfFinder {
	nf.defaultconf = cfgpath
	return nf
}

// CleanPaths removes all the paths from the stack
func (nf *NanoconfFinder) CleanPaths() *NanoconfFinder {
	nf.paths = nil
	return nf
}

// CleanNames removes all the names from the stack
func (nf *NanoconfFinder) CleanNames() *NanoconfFinder {
	nf.names = nil
	return nf
}

// AddName adds a name to the stack
func (nf *NanoconfFinder) AddName(name string) *NanoconfFinder {
	nf.names = append(nf.names, name)
	return nf
}

// AddPath adds a path to the stack
func (nf *NanoconfFinder) AddPath(cfgpath string) *NanoconfFinder {
	nf.paths = append(nf.paths, cfgpath)
	return nf
}
