function CertIsAboutToExpire(dateString) {
    return (new Date(dateString)).getTime() - Date.now() < 7 * 86400;
}
Vue.component('site-state', {
    props: ['state'],
    data: function () {
        let line1 = '';
        let line2 = '';

        return {
            line1: line1,
            line2: line2,
        }
    },
    template: '<div classs="alert-success">{{line1}}<br>{{line2}}</div>'
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
                        if (rec.SSLError !== null)
                            rec.Line2 = rec.SSLError;
                        else {
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