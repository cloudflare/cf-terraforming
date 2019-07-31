package cmd

import (
	"os"

	"text/template"

	cloudflare "github.com/cloudflare/cloudflare-go"
	"github.com/spf13/cobra"

	"github.com/sirupsen/logrus"
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
		log.Debug("Importing Account Member data")

		if accountID == "" {
			log.Error("'account' must be set.")
		}

		accountMembers, _, err := api.AccountMembers(accountID, cloudflare.PaginationOptions{
			Page:    1,
			PerPage: 1000,
		})

		if err != nil {
			log.Debug(err)
		}

		for _, r := range accountMembers {

			log.WithFields(logrus.Fields{
				"Account member ID": r.ID,
				"Status":            r.Status,
				"User ID":           r.User.ID,
				"User email":        r.User.Email,
			}).Debug("Processing account member")

			if tfstate {
				// TODO: Implement state dump
			} else {
				memberParse(r)
			}
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
