package actions

type ActionResult struct {
  Error interface{} `json:"error"`
  Response string `json:"response"`
}
