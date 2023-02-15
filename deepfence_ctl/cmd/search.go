package cmd

import (
	"context"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/deepfence/ThreatMapper/deepfence_ctl/http"
	"github.com/deepfence/ThreatMapper/deepfence_ctl/output"
	deepfence_server_client "github.com/deepfence/golang_deepfence_sdk/client"
	"github.com/deepfence/golang_deepfence_sdk/utils/log"
)

func type2field(t string) string {
	switch t {
	case "vulnerability":
		return "vulnerabilities_count"
	case "secret":
		return "secrets_count"
	case "compliance":
		return "compliances_count"
	case "cloudcompliance":
		return "cloudcompliances_count"
	case "malware":
		return "malwares_count"
	}
	return ""
}

var topCmd = &cobra.Command{
	Use:   "top",
	Short: "Graph Node search",
	Long:  `This subcommand fetch details for graph nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		search_type, _ := cmd.Flags().GetString("toptype")
		target_type, _ := cmd.Flags().GetString("type")

		search_fields, _ := cmd.Flags().GetString("fields")
		fields := []string{}
		if len(search_fields) != 0 {
			fields = strings.Split(search_fields, ",")
		}

		rank, _ := cmd.Flags().GetString("num")
		num, err := strconv.Atoi(rank)
		if err != nil {
			log.Fatal().Msgf("%v", err)
		}

		orderFilter := deepfence_server_client.ReportersOrderFilter{
			Limit:      int32(num),
			OrderField: type2field(search_type),
		}

		filters := deepfence_server_client.ReportersSearchFilter{
			InFieldFilter: fields,
			Filters: deepfence_server_client.ReportersFieldsFilters{
				ContainsFilter: deepfence_server_client.ReportersContainsFilter{},
				OrderFilter:    orderFilter,
			},
		}

		switch target_type {
		case "host":
			req := http.Client().LookupApi.SearchHosts(context.Background())
			req = req.ReportersSearchFilter(filters)
			res, rh, err := http.Client().LookupApi.SearchHostsExecute(req)
			if err != nil {
				log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			}
			output.Out(res)
		case "container":
			req := http.Client().LookupApi.SearchContainers(context.Background())
			req = req.ReportersSearchFilter(filters)
			res, rh, err := http.Client().LookupApi.SearchContainersExecute(req)
			if err != nil {
				log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			}
			output.Out(res)
		case "image":
			req := http.Client().LookupApi.SearchContainerImages(context.Background())
			req = req.ReportersSearchFilter(filters)
			res, rh, err := http.Client().LookupApi.SearchContainerImagesExecute(req)
			if err != nil {
				log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			}
			output.Out(res)
			//	case "vulnerability":
			//		req := http.Client().LookupApi.SearchContainerImages(context.Background())
			//		req = req.ReportersSearchFilter(filters)
			//		res, rh, err := http.Client().LookupApi.SearchContainerImagesExecute(req)
			//		if err != nil {
			//			log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			//		}
			//		output.Out(res)
			//	case "secret":
			//		req := http.Client().LookupApi.SearchContainerImages(context.Background())
			//		req = req.ReportersSearchFilter(filters)
			//		res, rh, err := http.Client().LookupApi.SearchContainerImagesExecute(req)
			//		if err != nil {
			//			log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			//		}
			//		output.Out(res)
			//	case "malware":
			//		req := http.Client().LookupApi.SearchContainerImages(context.Background())
			//		req = req.ReportersSearchFilter(filters)
			//		res, rh, err := http.Client().LookupApi.SearchContainerImagesExecute(req)
			//		if err != nil {
			//			log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			//		}
			//		output.Out(res)
			//	case "compliance":
			//		req := http.Client().LookupApi.SearchContainerImages(context.Background())
			//		req = req.ReportersSearchFilter(filters)
			//		res, rh, err := http.Client().LookupApi.SearchContainerImagesExecute(req)
			//		if err != nil {
			//			log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			//		}
			//		output.Out(res)
			//	case "cloud-compliance":
			//		req := http.Client().LookupApi.SearchContainerImages(context.Background())
			//		req = req.ReportersSearchFilter(filters)
			//		res, rh, err := http.Client().LookupApi.SearchContainerImagesExecute(req)
			//		if err != nil {
			//			log.Fatal().Msgf("Fail to execute: %v: %v", err, rh)
			//		}
			//		output.Out(res)
		default:
			log.Fatal().Msgf("Unsupported type:%s", target_type)
		}

	},
}

func init() {
	rootCmd.AddCommand(topCmd)

	topCmd.PersistentFlags().String("type", "", "host/container/process")
	topCmd.PersistentFlags().String("toptype", "", "host/container/process")
	topCmd.PersistentFlags().String("num", "", "Number of top")
	topCmd.PersistentFlags().String("fields", "", "CSV fields to search")
}
