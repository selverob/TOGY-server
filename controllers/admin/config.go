package admin

import (
	"models/action"
	"models/configuration/config"
	"models/configuration/timeConfig"
	"net/http"
	"strconv"
	"time"
	"util"
)

//ShowConfig handles showing the page in which user can see and edit
//the central configuration for clients.
func ShowConfig(c util.Context) {
	conf, err := config.Get(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	as, err := action.GetFor(&config.Config{}, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}

	a := prepareActions(as)

	util.RenderLayout("config.html", "Konfigurácia obrazoviek", struct {
		Conf     config.Config
		A        map[string][]time.Time
		ZeroTime time.Time
	}{conf, a, time.Date(0001, 01, 01, 00, 00, 00, 00, utc)}, c, "/static/js/jquery-ui-1.9.2.custom.min.js", "/static/js/timepicker-min.js", "/static/js/config.js")
}

//SetConfig handles saving the new configuration to Datastore.
func SetConfig(c util.Context) {
	var err error
	conf := new(config.Config)
	conf.StandardOn, err = time.Parse(config.ConfTimeFormat, c.R.FormValue("standardOn"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	conf.StandardOff, err = time.Parse(config.ConfTimeFormat, c.R.FormValue("standardOff"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	conf.OverrideState, err = strconv.Atoi(c.R.FormValue("overrideState"))
	if err != nil {
		util.Log500(err, c)
		return
	}

	conf.UpdateInterval, err = strconv.Atoi(c.R.FormValue("updateInterval"))
	if err != nil {
		util.Log500(err, c)
		return
	}

	if c.R.FormValue("weekends") == "true" {
		conf.Weekends = true
	}

	err = conf.Save(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	http.Redirect(c.W, c.R, "/admin/config", 303)
}

//TimeOverride renders the list of time overrides in Datastore.
func TimeOverride(c util.Context) {
	tcs, err := timeConfig.GetAll(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	util.RenderLayout("timeConfig.html", "Zoznam časových výnimiek", tcs, c)
}

//TimeOverrideEdit handles editing of existing time overrides.
//If it doesn't find id value in the path, it adds a new override.
func TimeOverrideEdit(c util.Context) {
	var tc *timeConfig.TimeConfig
	var err error
	if key := c.Vars["id"]; key == "" {
		tc = nil
	} else {
		tc, err = timeConfig.GetByKey(key, c.Ac)
		if err != nil {
			util.Log500(err, c)
			return
		}
	}
	util.RenderLayout("timeConfigEdit.html", "Úprava výnimky", tc, c, "/static/js/jquery-ui-1.9.2.custom.min.js", "/static/js/timepicker-min.js", "/static/js/editTC.js")

}

//TimeOverrideSubmit handles saving of time overrides into Datastore.
func TimeOverrideSubmit(c util.Context) {
	date, err := time.Parse(config.ConfDateFormat, c.R.FormValue("date"))
	if err != nil {
		util.Log500(err, c)
		return
	}

	on, err := time.Parse(config.ConfTimeFormat, c.R.FormValue("on"))
	if err != nil {
		util.Log500(err, c)
		return
	}

	off, err := time.Parse(config.ConfTimeFormat, c.R.FormValue("off"))
	if err != nil {
		util.Log500(err, c)
		return
	}
	tc := timeConfig.New(util.NormalizeDate(date, true), util.NormalizeTime(on, true), util.NormalizeTime(off, true))
	tc.Key = c.Vars["id"]
	err = tc.Save(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	http.Redirect(c.W, c.R, "/admin/config/timeOverride", 303)
}

//TimeOverrideDelete handles deleting of a time override.
func TimeOverrideDelete(c util.Context) {
	key := c.R.FormValue("key")
	tc, err := timeConfig.GetByKey(key, c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	err = tc.Delete(c.Ac)
	if err != nil {
		util.Log500(err, c)
		return
	}
	http.Redirect(c.W, c.R, "/admin/config/timeOverride", 303)
}
