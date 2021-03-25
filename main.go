package main

import (
	"github.com/valyala/fasthttp"
)

const (
	servicePort = `:8080`
)

var sampleItems = []int{2186603, 2654554, 2968009, 3723312, 4144387, 4588641, 4815118, 4861382, 5024530,
	5146236, 5613752, 5616296, 6170054, 6590837, 6671965, 6790848, 6796776, 7601411, 8139643, 8781920,
	9024198, 9322369, 9434636, 9434637, 9434638, 10004476, 10143840, 10219059, 10459579, 10667312,
	10854001, 10854002, 11088253, 11481196, 11857229, 12016099, 12016103, 12423144, 12507082, 12888079,
	12968083, 13227933, 13595361, 13615120, 13615122, 13807725, 13918691, 13944386, 13978470, 14136935,
	14324471, 14878237, 15025349, 15069435, 15154431, 15163808, 15435674, 15457224, 15556061, 15556062,
	16023989, 16379986, 16780324, 16889371, 17367533, 18029960, 18247707, 18362879, 18565187, 18622848}

var sampleUrls = []string{
	`http://catalog-backend-part3.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=3723312;`,
	`http://catalog-backend-part4.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=4144387;4588641;4815118;4861382;`,
	`http://catalog-backend-part5.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=5024530;5146236;5613752;5616296;`,
	`http://catalog-backend-part6.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=6170054;6590837;6671965;6790848;6796776;`,
	`http://catalog-backend-part7.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=7601411;`,
	`http://catalog-backend-part8.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=8139643;8781920;`,
	`http://catalog-backend-part9.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=9024198;9322369;9434636;9434637;9434638;`,
	`http://catalog-backend-part10.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=10004476;10143840;10219059;10459579;10667312;10854001;10854002;`,
	`http://catalog-backend-part11.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=11088253;11481196;11857229;`,
	`http://catalog-backend-part12.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=12016099;12016103;12423144;12507082;12888079;12968083;`,
	`http://catalog-backend-part13.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=13227933;13595361;13615120;13615122;13807725;13918691;13944386;13978470;`,
	`http://catalog-backend-part14.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=14136935;14324471;14878237;`,
	`http://catalog-backend-part15.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=15025349;15069435;15154431;15163808;15435674;15457224;15556061;15556062;`,
	`http://catalog-backend-part16.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=16023989;16379986;16780324;16889371;`,
	`http://catalog-backend-part17.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=17367533;`,
	`http://catalog-backend-part18.wbx-ru.svc.k8s.dataline/catalog?lang=ru&locale=ru&product=18029960;18247707;18362879;18565187;18622848;`,
}

func serviceHandler(ctx *fasthttp.RequestCtx) {
	w := defaultPool.acquire()
	w.run(16)
	w.release()
	ctx.WriteString(`Done!`)
}

func main() {
	fasthttp.ListenAndServe(servicePort, serviceHandler)
}
