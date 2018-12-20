package playground

import (
	"fmt"
	"io/ioutil"
	"margo.sh/mg"
	"net/http"
	"strings"
)

type Share struct{ mg.ReducerType }

func (r *Share) RCond(mx *mg.Ctx) bool { return mx.LangIs(mg.Go) }

func (r *Share) Reduce(mx *mg.Ctx) *mg.State {
	switch mx.Action.(type) {
	case mg.RunCmd:
		return mx.AddBuiltinCmds(mg.BuiltinCmd{
			Name: "GoPlayground.Share",
			Run: func(cx *mg.CmdCtx) *mg.State {
				go r.share(cx)
				return cx.State
			},
		})
	case mg.QueryUserCmds:
		return mx.AddUserCmds(mg.UserCmd{
			Title:   "Share file on play.golang.org",
			Name:    "GoPlayground.Share",
			Prompts: []string{"Type `share` to confirm sharing"},
		})
	}
	return mx.State
}

func (r *Share) share(cx *mg.CmdCtx) {
	defer cx.Output.Close()

	if l := cx.Prompts; len(l) != 1 || strings.ToLower(l[0]) != "share" {
		fmt.Fprintln(cx.Output, "Error: Please type `share` at the prompt.")
		return
	}

	body, err := cx.View.Open()
	if err != nil {
		fmt.Fprintln(cx.Output, "Error: cannot open view:", err)
	}
	defer body.Close()

	u := "https://play.golang.org"
	req, _ := http.NewRequest("POST", u+"/share", body)
	req.Header.Set("User-Agent", "GoPlayground.Share (margo.sh")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintln(cx.Output, "Error: request failed:", err)
		return
	}
	defer resp.Body.Close()

	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(cx.Output, "Error: cannot read response:", err)
	}
	if err != nil {
		fmt.Fprintln(cx.Output, "Error:", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Fprintln(cx.Output, "Error: unexpected http status:", resp.Status)
		return
	}
	fmt.Fprintf(cx.Output, "%s/p/%s\n", u, s)
}
