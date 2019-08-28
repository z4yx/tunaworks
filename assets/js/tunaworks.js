function CertIsAboutToExpire(dateString) {
    let epoch = (new Date(dateString)).getTime();
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
        secondRow: [],
        nodeInfo: [],
        websites: [],
    },
    methods: {
        buildCellContent: function(info, site, protocol) {
            let nodes = protocol == 4 ? site.Nodes4 : site.Nodes6;
            for (let node_id in nodes) {
                let rec = nodes[node_id];
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
                        rec.Line1 = "ERR";
                        if(rec.SSLError.length > 20){
                            let suffix = rec.SSLError.slice(-20);
                            if(suffix.indexOf(' ') > 0)
                                suffix = suffix.slice(suffix.indexOf(' '));
                            rec.Line2 = '...' + suffix;
                            rec.Details = rec.SSLError;
                        }else
                            rec.Line2 = rec.SSLError;
                    } else {
                        rec.Line1 = rec.StatusCode;
                        rec.Line2 = rec.ResponseTime + " ms";
                    }
                } else if (site.Url.startsWith("https:") && CertIsAboutToExpire(rec.SSLExpire)) {
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
                // console.debug(rec);
            }
            for (let node_id in info.NodeNames)
                if (!(node_id in nodes))
                    nodes[node_id] = {};
        },
        buildDispProp: function (info) {
            let secondRow = [];
            let nodeInfo = [];
            for (let node_id in info.NodeNames) {
                let proto = info.NodeInfo[node_id].Proto;
                let nProto = 0;
                if(proto & 1){
                    secondRow.push("IPv4");
                    nProto++;
                }
                if(proto & 2){
                    secondRow.push("IPv6");
                    nProto++;
                }
                if(nProto == 0)
                    continue;
                nodeInfo.push([node_id, info.NodeNames[node_id], info.NodeInfo[node_id].Heartbeat, nProto]);
            }
            this.nodeInfo = nodeInfo;
            this.secondRow = secondRow;
            let websites = [];
            info.Websites.forEach((site) => {
                this.buildCellContent(info, site, 4);
                this.buildCellContent(info, site, 6);
                let records = [];
                for (let pair of nodeInfo) { // order matters
                    let node_id = pair[0];
                    let proto = info.NodeInfo[node_id].Proto;
                    if(proto & 1)
                        records.push(site.Nodes4[node_id]);
                    if(proto & 2)
                        records.push(site.Nodes6[node_id]);
                }
                websites.push({
                    Url: site.Url,
                    Records: records,
                });
            });

            this.websites = websites;
        },
        loadLatest: function () {
            this.$http.get('monitor/latest').then((resp) => {
                let body = resp.body;
                this.buildDispProp(body);
            }, () => {

            });
        }
    },
    filters: {
        moment: function(date){
            // console.debug(date);
            if(date.startsWith("1970-"))
                return 'never';
            return moment(date).fromNow();
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
                let nodes = [];
                if(site.Url.startsWith("http:"))
                    return;
                for (let node in site.Nodes4)
                    nodes.push(site.Nodes4[node]);
                for (let node in site.Nodes6)
                    nodes.push(site.Nodes6[node]);
                for (let item of nodes) {
                    if(rec === null) {
                        rec = item;
                        continue;
                    }
                    // We're interested in x509 errors only
                    if(item.SSLError !== null && !item.SSLError.startsWith('x509'))
                        continue;
                    if(rec.SSLError !== null && item.SSLError === null){
                        rec = item;
                        continue;
                    }
                    if(rec.SSLError === null && item.SSLError !== null){
                        continue;
                    }
                    let updated = (new Date(item.Updated)).getTime();
                    if(updated > latestUpdate) {
                        latestUpdate = updated;
                        rec = item;
                    }
                }
                if(rec === null)
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
                    item.Icon['fa-times-circle'] = true;
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
