package cmd

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/oworkshop/ow_wiki/app/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// NewCommand returns cobra.Command to run
func NewCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "ow_wiki",
		Short: "ow_wiki: a simple wiki system",
		Run: func(cmd *cobra.Command, args []string) {
			if ok := InitConfig(); !ok {
				return
			}
			startServer()
		},
	}

	cmds.ResetFlags()

	pflag.StringVarP(&MainConfigFile, "config", "c", DefaultConfigFile, "config file")
	pflag.IntVarP(&ListenPort, "port", "p", 0, "listen port (highest priority)")

	cmds.AddCommand(newCmdConfig(out))

	return cmds
}

func startServer() error {
	router := mux.NewRouter()
	router.PathPrefix("/css/").Handler(http.StripPrefix("/css/",
		http.FileServer(http.Dir(filepath.Join(utils.Basepath(), "template/css/")))))
	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/",
		http.FileServer(http.Dir(filepath.Join(utils.Basepath(), "template/js/")))))
	router.PathPrefix("/").Handler(
		negroni.New(
			negroni.HandlerFunc(checkAccessAllowed),
			negroni.HandlerFunc(NewHandler),
		))

	server := negroni.New(negroni.NewRecovery())
	server.UseHandler(router)

	return http.ListenAndServe("0.0.0.0:"+mainConfig.ListenPort, server)
}
