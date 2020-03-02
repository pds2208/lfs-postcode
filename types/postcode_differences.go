package types

type PostCodeDifferences struct {
    CaseNo   int    `db:"case_no" json:"caseNo"`
    PostCode string `db:"post_code" json:"postCode"`
}
