# GOS 
ALL in One tool
```
=======     ===     =======
=          =   =   = 
=         =     =  =
=  =====  =     =   ======
=      =  =     =         =
=     =    =   =          =
 =======    ===    =======
```

Features:
* Race Condition
* Email Phishing
* Brute-Force SSH
* Brute-Force LDAP
* TCP Forward Port

Support:
* [x] Windows
* [x] Linux

## Usage
```
Usage: gos <option> <agrs>

Options:
        race    Testing for Race Condition
        mail    Testing for Email Phishing
        ssh     Brute-Force SSH
        ldap    Brute-Force LDAP
        tcp     TCP Forward Port
EXAMPLES:
        gos race -r req.txt -p http://burp:8080 -w 10 -ssl
        gos mail --from support@gmail.com --to victim@company.org --attach passwd.txt --attach-name "../etc/passwd"
        gos ssh -H hosts.txt -P pass.txt -u admin
        gos ldap -s 10.10.10.1 -P pass.txt -u admin -d example.com
        gos tcp -s 192.168.0.10:3389 -d 10.10.10.1:3389

Use "gos <option> --h" for more information about a option.
```
## Example
Race Condition: [req.txt](https://raw.githubusercontent.com/namcuongq/gos/main/example/req.txt) . You can copy it from burpsuite


## Download

[gos.exe](https://github.com/namcuongq/gos/releases)

## TODO

* [ ] Brute-Force RDP
* [ ] HTTP server - Directory indexing
