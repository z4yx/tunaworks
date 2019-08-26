function CertIsAboutToExpire(dateString) {
    let epoch = (new Date(dateString)).getTime();
    if (epoch === 0) return false;
    return epoch - Date.now() < 7 * 86400 * 1000;
}
Vue.component('top-nav-bar', {
    props: ['current'],
    data: function () {
        return {
            links: [['./', 'Overall'], ['./ssl', 'SSL']],
        }
    },
    template: '<nav class="navbar navbar-expand-lg navbar-light" style="background-color: #e3f2fd;">\
      <a class="navbar-brand" href="#">TUNA.works</a>\
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent"\
        aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">\
        <span class="navbar-toggler-icon"></span>\
      </button>\
      <div class="collapse navbar-collapse" id="navbarSupportedContent">\
        <ul class="navbar-nav mr-auto">\
          <li v-for="link in links" v-bind:class="{ \'nav-item\': true, active: current===link[1] }">\
            <a class="nav-link" v-bind:href="link[0]">{{link[1]}}</a>\
          </li>\
        </ul>\
      </div>\
    </nav>',
})
var overall = document.getElementById('overall') && new Vue({
    el: '#overall',
    data: {
        latestInfo: {
            NodeNames: {},
            Websites: [],
        }
    },
    methods: {
        insertDispProp: function (info) {
            info.Websites.forEach((site) => {
                for (let node in site.Nodes) {
                    let rec = site.Nodes[node];
                    rec.ClassObj = {
                        "cell-site-state": true
                    };
                    rec.Icon = {
                        "fas": true
                    };
                    if (rec.SSLError !== null || rec.StatusCode >= 400) {
                        rec.ClassObj["alert-danger"] = true;
                        rec.Icon["fa-times"] = true;
                        if (rec.SSLError !== null) {
                            rec.Line1 = "SSL";
                            rec.Line2 = rec.SSLError;
                        } else {
                            rec.Line1 = rec.StatusCode;
                            rec.Line2 = rec.ResponseTime + " ms";
                        }
                    } else if (CertIsAboutToExpire(rec.SSLExpire)) {
                        rec.ClassObj["alert-warning"] = true;
                        rec.Icon["fa-exclamation-triangle"] = true;
                        rec.Line1 = rec.StatusCode;
                        rec.Line2 = 'Cert will expire on ' + (new Date(rec.SSLExpire)).toLocaleDateString();
                    } else {
                        rec.ClassObj["alert-success"] = true;
                        rec.Icon["fa-check"] = true;
                        rec.Line1 = rec.StatusCode;
                        rec.Line2 = rec.ResponseTime + " ms";
                    }
                    console.debug(rec);
                }
                for (let node in info.NodeNames)
                    if (!(node in site.Nodes))
                        site.Nodes[node] = {};
            });
        },
        loadLatest: function () {
            this.$http.get('monitor/latest').then((resp) => {
                let body = resp.body;
                this.insertDispProp(body);
                this.latestInfo = body;
            }, () => {

            });
        }
    },
    created: function () {
        this.loadLatest();
    },
})
var ssl = document.getElementById('ssl') && new Vue({
    el: '#ssl',
    data: {
        sslInfo: []
    },
    methods: {
        genSSLInfo: function (info) {
            let sslInfo = [];
            info.Websites.forEach((site) => {
                let latestUpdate = 0;
                let rec = null;
                for (let node in site.Nodes) {
                    let item = site.Nodes[node];
                    let updated = (new Date(item.Updated)).getTime();
                    if(updated > latestUpdate) {
                        latestUpdate = updated;
                        rec = item;
                    }
                }
                if(rec === null || (rec.SSLError === null && rec.SSLExpire.startsWith('1970')))
                    return;
                let item = {
                    ClassObj: {},
                    Icon: {fas: true},
                    Url: site.Url,
                    Updated: rec.Updated,
                    Prober: rec.Name,
                    SortKey: 0,
                };
                if(rec.SSLError !== null) {
                    item.Icon['fa-times'] = true;
                    item.ClassObj['table-danger'] = true;
                    item.Expiry = rec.SSLError;
                } else if(CertIsAboutToExpire(rec.SSLExpire)){
                    item.Icon['fa-exclamation-circle'] = true;
                    item.ClassObj['table-warning'] = true;
                    item.Expiry = rec.SSLExpire;
                    item.SortKey = (new Date(rec.SSLExpire)).getTime();
                } else {
                    item.Icon['fa-calendar-alt'] = true;
                    item.Expiry = rec.SSLExpire;
                    item.SortKey = (new Date(rec.SSLExpire)).getTime();
                }
                sslInfo.push(item);
                console.debug(item);
            });
            sslInfo.sort((a,b)=>{
                return a.SortKey - b.SortKey;
            });
            this.sslInfo = sslInfo;
        },
        loadLatest: function () {
            this.$http.get('monitor/latest').then((resp) => {
                let body = resp.body;
                this.genSSLInfo(body);
            }, () => {

            });
        }
    },
    created: function () {
        this.loadLatest();
    },
    filters: {
        ISO8601: function(val) {
            let d= new Date(val);
            if(isNaN(d.getTime()))
                return val;
            else
                return d.toLocaleString();
        }
    },
});