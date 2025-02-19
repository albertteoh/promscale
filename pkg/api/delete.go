// This file and its contents are licensed under the Apache License 2.0.
// Please see the included NOTICE for copyright information and
// LICENSE for a copy of the license.

package api

import (
	"fmt"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/prometheus/prometheus/promql/parser"
	"github.com/timescale/promscale/pkg/log"
	"github.com/timescale/promscale/pkg/pgclient"
	"github.com/timescale/promscale/pkg/pgmodel/common/errors"
	deletePkg "github.com/timescale/promscale/pkg/pgmodel/delete"
	"github.com/timescale/promscale/pkg/pgmodel/model"
)

func Delete(conf *Config, client *pgclient.Client) http.Handler {
	hf := corsWrapper(conf, deleteHandler(conf, client))
	return gziphandler.GzipHandler(hf)
}

func deleteHandler(config *Config, client *pgclient.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.ReadOnly {
			respondError(w, http.StatusForbidden, fmt.Errorf("read-only connector cannot perform deletion"), "operation_not_permitted")
			return
		}
		if !config.AdminAPIEnabled {
			respondError(w, http.StatusForbidden, fmt.Errorf("deletion of series requires admin permissions. Use -web-enable-admin-api flag to allow deletion operations"), "operation_not_permitted")
			return
		}
		var (
			totalRowsDeleted int
			metricsTouched   []string
			seriesDeleted    []model.SeriesID
		)
		if err := r.ParseForm(); err != nil {
			respondError(w, http.StatusBadRequest, err, "bad_data")
			return
		}
		if len(r.Form["match[]"]) == 0 {
			respondError(w, http.StatusBadRequest, fmt.Errorf("no match[] parameter provided"), "bad_data")
			return
		}
		start, err := parseTimeParam(r, "start", model.MinTime)
		if err != nil {
			log.Info("msg", "Query bad request:"+err.Error())
			respondError(w, http.StatusBadRequest, err, "bad_data")
			return
		}
		end, err := parseTimeParam(r, "end", model.MaxTime)
		if err != nil {
			log.Info("msg", "Query bad request:"+err.Error())
			respondError(w, http.StatusBadRequest, err, "bad_data")
			return
		}
		if start != model.MinTime || end != model.MaxTime {
			log.Warn("msg", "Time based series deletion is unsupported.")
			respondError(w, http.StatusBadRequest, errors.ErrTimeBasedDeletion, "bad_data")
			return
		}
		for _, s := range r.Form["match[]"] {
			matchers, err := parser.ParseMetricSelector(s)
			if err != nil {
				respondError(w, http.StatusBadRequest, err, "bad_data")
				return
			}
			if client == nil {
				continue
			}
			pgDelete := deletePkg.PgDelete{Conn: client.ReadOnlyConnection()}
			touchedMetrics, deletedSeriesIDs, rowsDeleted, err := pgDelete.DeleteSeries(r.Context(), matchers, start, end)
			if err != nil {
				respondErrorWithMessage(w, http.StatusInternalServerError, err, "deleting_series",
					fmt.Sprintf("partial delete: deleted %v series IDs from %v metrics, affecting %d rows in total.",
						distinctValues(seriesDeleted),
						distinctValues(metricsTouched),
						totalRowsDeleted,
					),
				)
				return
			}
			metricsTouched = append(metricsTouched, touchedMetrics...)
			seriesDeleted = append(seriesDeleted, deletedSeriesIDs...)
			totalRowsDeleted += rowsDeleted
		}
		respond(w, http.StatusOK,
			fmt.Sprintf("deleted %v series IDs from %v metrics, affecting %d rows in total.",
				distinctValues(seriesDeleted),
				distinctValues(metricsTouched),
				totalRowsDeleted,
			),
		)
	}
}

func distinctValues(slice interface{}) []string {
	temp := make(map[string]struct{})
	switch elem := slice.(type) {
	case []string:
		for _, element := range elem {
			if _, ok := temp[element]; !ok {
				temp[element] = struct{}{}
			}
		}
	case []model.SeriesID:
		for _, element := range elem {
			str := element.String()
			if _, ok := temp[str]; !ok {
				temp[str] = struct{}{}
			}
		}
	}
	keys := make([]string, 0, len(temp))
	for k := range temp {
		keys = append(keys, k)
	}
	return keys
}
