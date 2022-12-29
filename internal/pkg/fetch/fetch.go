package fetch

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/format"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/logger"
	"gitlab.com/EysteinnSig/stackmap-mapserver/internal/pkg/types"
)

func FetchAllProducts(outdir string, apihost string, sqldata types.SQLData) error {
	prods := types.UniqueProducts{}
	err := GetData(apihost+"/api/v1/products", &prods)
	if err != nil {
		return err
	}

	for _, prod := range prods.Products {
		logger.GetLogger().Debugf("Fetching product: %s\n", prod)
		err := FetchProduct(prod, outdir, apihost, sqldata)
		if err != nil {
			return err
		}
	}

	prodstr := ""
	for _, prod := range prods.Products {
		prodstr += fmt.Sprintf("INCLUDE \"product_%s.map\"\n", prod)

	}
	err = os.WriteFile(filepath.Join(outdir, "products.map"), []byte(prodstr), 0644)
	if err != nil {
		return err
	}

	return nil
}

func FetchProduct(product string, outdir string, apihost string, sqldata types.SQLData) error {

	//product := "hrit-ash"
	//pt := ProductTimes{}
	pt := types.MapLayerData{SQLData: sqldata}
	//GetData("http://localhost:3000/api/v1/times?product="+product, &pt)
	err := GetData(apihost+"/api/v1/times?product="+product, &pt)
	if err != nil {
		return err
	}
	sort.SliceStable(pt.Times, func(i, j int) bool {
		return pt.Times[i].Before(pt.Times[j])
	})
	pt.StartRange = pt.Times[0]
	pt.EndRange = pt.Times[len(pt.Times)-1]
	pt.DefaultTime = pt.EndRange

	/*pt.SQLDB = sqldata.SQLDB
	pt.SQLHost = sqldata.SQLHost
	pt.SQLPass = sqldata.SQLPass
	pt.SQLUser = sqldata.SQLUser*/

	//fmt.Println(pt.AllTimesString())
	//txt, err := Replace(product, timesstr, mintime, maxtime, mintime)
	bufr, err := format.GetMapfile(&pt)
	if err != nil {
		/*fmt.Println(err)
		os.Exit(1)*/
		return err
	}
	err = os.WriteFile(filepath.Join(outdir, "product_"+product+".map"), bufr, 0644)
	if err != nil {
		return err
	}
	return nil
}
