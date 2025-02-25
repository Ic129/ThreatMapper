package reports

import (
	"context"
	"os"

	"github.com/deepfence/ThreatMapper/deepfence_utils/log"
	"github.com/deepfence/ThreatMapper/deepfence_utils/utils"
	"github.com/xuri/excelize/v2"
)

var (
	vulnerabilityHeader = map[string]string{
		"A1": "@timestamp",
		"B1": "cve_attack_vector",
		"C1": "cve_caused_by_package",
		"D1": "cve_container_image",
		"E1": "scan_id",
		"F1": "cve_container_image_id",
		"G1": "cve_cvss_score",
		"H1": "cve_description",
		"I1": "cve_fixed_in",
		"J1": "cve_id",
		"K1": "cve_link",
		"L1": "cve_severity",
		"M1": "cve_overall_score",
		"N1": "cve_type",
		"O1": "host",
		"P1": "host_name",
		"Q1": "masked",
	}
	secretHeader = map[string]string{
		"A1": "Filename",
		"B1": "Content",
		"C1": "Name",
		"D1": "Rule",
		"E1": "Severity",
		"F1": "Node Name",
		"G1": "Container Name",
		"H1": "Kubernetes Cluster Name",
		"I1": "Signature",
	}
	malwareHeader = map[string]string{
		"A1": "Rule Name",
		"B1": "Severity",
		"C1": "Meta",
		"D1": "Meta Rules",
		"E1": "File Severity Score",
		"F1": "File Severity",
		"G1": "Summary",
		"H1": "Node Name",
		"I1": "Container Name",
		"J1": "Kubernetes Cluster Name",
		"K1": "NodeType",
	}
	complianceHeader = map[string]string{
		"A1": "@timestamp",
		"B1": "compliance_check_type",
		"C1": "count",
		"D1": "doc_id",
		"E1": "host",
		"F1": "host_name",
		"G1": "masked",
		"H1": "node_id",
		"I1": "node_name",
		"J1": "node_type",
		"K1": "status",
		"L1": "test_category",
		"M1": "test_desc",
		"N1": "test_info",
		"O1": "test_number",
	}
)

func generateXLSX(ctx context.Context, params utils.ReportParams) (string, error) {

	var (
		xlsxFile string
		err      error
	)

	switch params.Filters.ScanType {
	case VULNERABILITY:
		xlsxFile, err = vulnerabilityXLSX(ctx, params)
	case SECRET:
		xlsxFile, err = secretXLSX(ctx, params)
	case MALWARE:
		xlsxFile, err = malwareXLSX(ctx, params)
	case COMPLIANCE:
		xlsxFile, err = complianceXLSX(ctx, params)
	case CLOUD_COMPLIANCE:
		xlsxFile, err = cloudComplianceXLSX(ctx, params)
	default:
		return "", ErrUnknownScanType
	}

	if err != nil {
		return "", err
	}

	return xlsxFile, nil
}

func xlsxSave(xlsx *excelize.File, params utils.ReportParams) (string, error) {
	// create a temp file to hold xlsx report
	temp, err := os.CreateTemp("", "report-*-"+reportFileName(params))
	if err != nil {
		return "", err
	}
	defer temp.Close()

	// save spreadsheet by the given path.
	if err := xlsx.SaveAs(temp.Name()); err != nil {
		log.Error().Err(err).Msg("failed to save xlsx file")
		return "", err
	}
	return temp.Name(), nil
}

func xlsxSetHeader(xlsx *excelize.File, sheet string, headers map[string]string) {
	for k, v := range headers {
		xlsx.SetCellValue(sheet, k, v)
	}
}

func vulnerabilityXLSX(ctx context.Context, params utils.ReportParams) (string, error) {
	data, err := getVulnerabilityData(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to get vulnerabilities info")
		return "", err
	}

	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()

	xlsxSetHeader(xlsx, "Sheet1", vulnerabilityHeader)

	offset := 0
	for _, nodeScanData := range data.NodeWiseData.ScanData {
		for i, v := range nodeScanData.ScanResults {
			cellName, err := excelize.CoordinatesToCellName(1, offset+i+2)
			if err != nil {
				log.Error().Err(err).Msg("error generating cell name")
			}
			value := []interface{}{
				nodeScanData.ScanInfo.UpdatedAt,
				v.Cve_attack_vector,
				v.Cve_caused_by_package,
				nodeScanData.ScanInfo.NodeName,
				nodeScanData.ScanInfo.ScanID,
				nodeScanData.ScanInfo.NodeID,
				v.Cve_cvss_score,
				v.Cve_description,
				v.Cve_fixed_in,
				v.Cve_id,
				v.Cve_link,
				v.Cve_severity,
				v.Cve_overall_score,
				v.Cve_type,
				nodeScanData.ScanInfo.HostName,
				nodeScanData.ScanInfo.HostName,
				v.Masked,
			}
			xlsx.SetSheetRow("Sheet1", cellName, &value)
		}
		offset = offset + len(nodeScanData.ScanResults)
	}

	return xlsxSave(xlsx, params)
}

func secretXLSX(ctx context.Context, params utils.ReportParams) (string, error) {
	data, err := getSecretData(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to get secrets info")
		return "", err
	}

	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()

	xlsxSetHeader(xlsx, "Sheet1", secretHeader)

	offset := 0
	for _, nodeScanData := range data.NodeWiseData.ScanData {
		for i, s := range nodeScanData.ScanResults {
			cellName, err := excelize.CoordinatesToCellName(1, offset+i+2)
			if err != nil {
				log.Error().Err(err).Msg("error generating cell name")
			}
			value := []interface{}{
				s.FullFilename,
				s.MatchedContent,
				s.Name,
				s.RuleID,
				s.Level,
				nodeScanData.ScanInfo.NodeName,
				nodeScanData.ScanInfo.ContainerName,
				nodeScanData.ScanInfo.KubernetesClusterName,
				s.SignatureToMatch,
			}
			xlsx.SetSheetRow("Sheet1", cellName, &value)
		}
		offset = offset + len(nodeScanData.ScanResults)
	}

	return xlsxSave(xlsx, params)
}

func malwareXLSX(ctx context.Context, params utils.ReportParams) (string, error) {
	data, err := getMalwareData(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to get malwares info")
		return "", err
	}

	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()

	xlsxSetHeader(xlsx, "Sheet1", malwareHeader)

	offset := 0
	for _, nodeScanData := range data.NodeWiseData.ScanData {
		for i, m := range nodeScanData.ScanResults {
			cellName, err := excelize.CoordinatesToCellName(1, offset+i+2)
			if err != nil {
				log.Error().Err(err).Msg("error generating cell name")
			}
			value := []interface{}{
				m.RuleName,
				m.SeverityScore,
				"",
				m.RuleID,
				m.FileSevScore,
				m.FileSeverity,
				m.Summary,
				nodeScanData.ScanInfo.NodeName,
				nodeScanData.ScanInfo.ContainerName,
				nodeScanData.ScanInfo.KubernetesClusterName,
				nodeScanData.ScanInfo.NodeType,
			}
			xlsx.SetSheetRow("Sheet1", cellName, &value)
		}
		offset = offset + len(nodeScanData.ScanResults)
	}

	return xlsxSave(xlsx, params)
}

func complianceXLSX(ctx context.Context, params utils.ReportParams) (string, error) {
	data, err := getComplianceData(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to get compliance info")
		return "", err
	}

	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()

	xlsxSetHeader(xlsx, "Sheet1", complianceHeader)

	offset := 0
	for _, nodeScanData := range data.NodeWiseData.ScanData {
		for i, c := range nodeScanData.ScanResults {
			cellName, err := excelize.CoordinatesToCellName(1, offset+i+2)
			if err != nil {
				log.Error().Err(err).Msg("error generating cell name")
			}
			value := []interface{}{
				nodeScanData.ScanInfo.UpdatedAt,
				c.ComplianceCheckType,
				"",
				"",
				nodeScanData.ScanInfo.HostName,
				nodeScanData.ScanInfo.HostName,
				c.Masked,
				c.ComplianceNodeId,
				nodeScanData.ScanInfo.NodeName,
				c.ComplianceNodeType,
				c.Status,
				c.TestCategory,
				c.TestDesc,
				c.TestInfo,
				c.TestNumber,
			}
			xlsx.SetSheetRow("Sheet1", cellName, &value)
		}
		offset = offset + len(nodeScanData.ScanResults)
	}

	return xlsxSave(xlsx, params)
}

func cloudComplianceXLSX(ctx context.Context, params utils.ReportParams) (string, error) {
	data, err := getCloudComplianceData(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to get cloud compliance info")
		return "", err
	}

	xlsx := excelize.NewFile()
	defer func() {
		if err := xlsx.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close file")
		}
	}()

	xlsxSetHeader(xlsx, "Sheet1", complianceHeader)

	for _, data := range data.NodeWiseData.ScanData {
		for i, c := range data.ScanResults {
			cellName, err := excelize.CoordinatesToCellName(1, i+2)
			if err != nil {
				log.Error().Err(err).Msg("error generating cell name")
			}
			value := []interface{}{
				data.ScanInfo.UpdatedAt,
				c.ComplianceCheckType,
				"",
				"",
				data.ScanInfo.HostName,
				data.ScanInfo.HostName,
				c.Masked,
				c.NodeID,
				data.ScanInfo.NodeName,
				c.ComplianceCheckType,
				c.Status,
				c.Type,
				c.Description,
				c.Title,
				c.ControlID,
			}
			xlsx.SetSheetRow("Sheet1", cellName, &value)
		}
	}

	return xlsxSave(xlsx, params)
}
