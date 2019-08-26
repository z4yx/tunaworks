function CertIsAboutToExpire(dateString) {
    return (new Date(dateString)).getTime() - Date.now() < 7 * 86400 * 1000;
}
Vue.component('top-nav-bar', {
    data: function () {
        return {
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
          <li class="nav-item active">\
            <a class="nav-link" href="./">Overall <span class="sr-only">(current)</span></a>\
          </li>\
          <li class="nav-item">\
            <a class="nav-link" href="./ssl">SSL</a>\
          </li>\
        </ul>\
      </div>\
    </nav>'
})
var overall = new Vue({
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
                        rec.Line2 = 'Cert expires on ' + (new Date(rec.SSLExpire)).toLocaleDateString();
                    } else {
                        rec.ClassObj["alert-success"] = true;
                        rec.Icon["fa-check"] = true;
                        rec.Line1 = rec.StatusCode;
                        rec.Line2 = rec.ResponseTime + " ms";
                    }
                    console.log(rec);
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