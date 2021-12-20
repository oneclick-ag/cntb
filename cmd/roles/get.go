package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"contabo.com/cli/cntb/client"
	contaboCmd "contabo.com/cli/cntb/cmd"
	"contabo.com/cli/cntb/cmd/util"
	"contabo.com/cli/cntb/outputFormatter"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var roleGetCmd = &cobra.Command{
	Use:   "role [type] [roleId]",
	Short: "Info about a specific role",
	Long:  `Retrieves information about one role identified by type and id. available permission types are resourcePermission and apiPermission`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, httpResp, err := client.ApiClient().RolesApi.RetrieveRole(context.Background(), permissionType, roleId).
			XRequestId(uuid.NewV4().String()).Execute()

		util.HandleErrors(err, httpResp, "while retrieving role")

		responseJson, _ := json.Marshal(resp.Data)

		configFormatter = outputFormatter.FormatterConfig{
			Filter:     []string{"roleId", "name"},
			WideFilter: []string{"roleId", "name"},
			JsonPath:   contaboCmd.OutputFormatDetails,
		}

		util.HandleResponse(responseJson, configFormatter)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			cmd.Help()
			log.Fatal("Please specify type and ")
		}
		permissionType = args[0]
		if permissionType != "apiPermission" && permissionType != "resourcePermission" {
			cmd.Help()
			log.Fatal("Permission type can only be on of the following either apiPermission or resourcePermission")
		}

		permissionId64, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatal(fmt.Sprintf("Specified permissionId64 %v is not valid", args[0]))
		}
		roleId = permissionId64
		contaboCmd.ValidateOutputFormat()
		return nil
	},
}

func init() {
	contaboCmd.GetCmd.AddCommand(roleGetCmd)
}
