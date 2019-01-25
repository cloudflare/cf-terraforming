package cmd

import (
	"fmt"
	"log"
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"
)

const accountMemberTemplate = `
resource "cloudflare_account_member" "{{.Member.ID}}" {
    email_address = "{{.Member.User.Email}}"
    role_ids = [{{range .Member.Roles}}"{{.ID}}",{{end}}]
}
`

func init() {
	rootCmd.AddCommand(accountMemberCmd)
}

var accountMemberCmd = &cobra.Command{
	Use:   "account_member",
	Short: "Import Account Member data into Terraform",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Importing Account Member data")

		if accountID == "" {
			fmt.Println("'account' must be set.")
			os.Exit(1)
		}

		accountMembers, _, err := api.AccountMembers(accountID, cloudflare.PaginationOptions{
			Page:    1,
			PerPage: 1000,
		})

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, r := range accountMembers {
			log.Printf("[DEBUG] Account Member ID %s, Status %s, User.ID %s, User.Email %s\n", r.ID, r.Status, r.User.ID, r.User.Email)
			memberParse(r)
		}
	},
}

func memberParse(member cloudflare.AccountMember) {
	tmpl := template.Must(template.New("script").Funcs(templateFuncMap).Parse(accountMemberTemplate))
	tmpl.Execute(os.Stdout,
		struct {
			Member cloudflare.AccountMember
		}{
			Member: member,
		})
}
