# indicators-of-compromise

This is the repo that builds IOC list(s) that are used by the detector to flag IOCs.

## Sources

We depend on outside sources for the majority of our IOCs. We will add to this list as we go along

1. Shallalist -- a collection of over a million IOCs that are meant to be used for firewalls, etc. Lots of them are invalid, and the list hasn't been updated since at least the middle of 2017. [[1]](http://www.shallalist.de/)
2. Dan's Tor Node List -- a list of known Tor nodes that is updated every half hour at most. [[2]](http://dan.me/)
3. Botherder's Targeted Threats list -- a list of known threats targetting civil society, including activist and journalists. It's maintained by someone working with Amnesty International. [[3]](https://github.com/botherder/)
4. Trust Networks -- IOCs that have been collected by us. [[4]](https://www.trustnetworks.com/)

## Categories

We categorise the IOCs for use in the visualisations, including Kibana and the various Risk Graphs.

    aggressive  
    anonvpn  
    compromised.creds  
    covert.dns-tunnel  
    drugs  
    dynamic  
    gamble  
    hacking  
    insecure  
    location.unexpected  
    lost.equipment  
    ntp.private  
    oppression.rights  
    physical.security.breach  
    policy.violation  
    porn  
    port.scan  
    redirector  
    sex.lingerie  
    spyware  
    test  
    tor.entry  
    tor.exit  
    violence  
    warez  

### Add a risk category

The risk categories are stored in gaffer, and that's where we get them from. To update the risk categories, edit the script `load-riskCategories.py` and attach to one of the gaffer pods using port forwarding:

```$ kubectl port-forward gaffer-[version]-[id] 8080:8080```

and then run:

```$ GAFFER=http://localhost:8080 ./load-riskCategories.py -w```

### Generate included files

The script `get-riskCategories.py` will generate the contents of the files `trustnetworks/web/include.js` and `trustnetworks/analytics-comms-trust/src/analytics/risk-params.go`.

Whenever you add a risk category and have updated it using `load-riskCategories.py` on all of the gaffers that need it, generate the included files and update `web` and `analytics-comms-trust` and deploy in the analytics clusters.
