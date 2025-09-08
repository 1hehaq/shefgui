package main

const windowUI = `<?xml version="1.0" encoding="UTF-8"?>
<interface>
  <menu id="main_menu">
    <section>
      <item>
        <attribute name="label">Preferences</attribute>
        <attribute name="action">app.preferences</attribute>
      </item>
    </section>
    <section>
      <item>
        <attribute name="label">Help</attribute>
        <attribute name="action">app.help</attribute>
      </item>
      <item>
        <attribute name="label">About</attribute>
        <attribute name="action">app.about</attribute>
      </item>
    </section>
  </menu>
  <requires lib="gtk" version="4.0"/>
  <object class="GtkApplicationWindow" id="window">
    <property name="title">shef</property>
    <property name="default-width">800</property>
    <property name="default-height">600</property>
    <child type="titlebar">
      <object class="GtkHeaderBar" id="header_bar">
        <property name="title-widget">
          <object class="GtkLabel">
            <property name="label">shef</property>
          </object>
        </property>
        <child type="start">
          <object class="GtkMenuButton" id="main_menu_button">
            <property name="icon-name">open-menu-symbolic</property>
            <property name="menu-model">main_menu</property>
          </object>
        </child>
        <child type="end">
          <object class="GtkButton" id="title_search_button">
            <property name="icon-name">system-search-symbolic</property>
            <property name="tooltip-text">Search in results</property>
            <style>
              <class name="flat"/>
            </style>
          </object>
        </child>
      </object>
    </child>
    <child>
      <object class="GtkBox" id="main_box">
        <property name="orientation">vertical</property>
        <property name="spacing">10</property>
        <property name="margin-top">20</property>
        <property name="margin-bottom">20</property>
        <property name="margin-start">20</property>
        <property name="margin-end">20</property>
        <child>
          <object class="GtkBox" id="search_box">
            <property name="orientation">horizontal</property>
            <property name="spacing">10</property>
            <child>
              <object class="GtkEntry" id="query_entry">
                <property name="placeholder-text">search query</property>
                <property name="hexpand">true</property>
                <property name="width-chars">30</property>
                <property name="max-width-chars">80</property>
              </object>
            </child>
            <child>
              <object class="GtkComboBoxText" id="facet_combo">
                <property name="active">0</property>
                <property name="hexpand">false</property>
              </object>
            </child>
            <child>
              <object class="GtkButton" id="search_button">
                <property name="label">Search</property>
                <style>
                  <class name="suggested-action"/>
                </style>
              </object>
            </child>
            <child>
              <object class="GtkSpinner" id="spinner">
                <property name="visible">false</property>
              </object>
            </child>
          </object>
        </child>

        <child>
          <object class="GtkBox" id="results_container">
            <property name="orientation">vertical</property>
            <property name="spacing">10</property>
            <property name="vexpand">true</property>
            <child>
              <object class="GtkBox" id="results_header_box">
                <property name="orientation">horizontal</property>
                <property name="spacing">6</property>
                <child>
                  <object class="GtkLabel" id="results_count">
                    <property name="label">0</property>
                    <property name="halign">start</property>
                    <style>
                      <class name="count-label"/>
                    </style>
                  </object>
                </child>
                <child>
                  <object class="GtkBox" id="spacer_box">
                    <property name="hexpand">true</property>
                  </object>
                </child>
                <child>
                  <object class="GtkBox" id="results_toolbar">
                    <property name="orientation">horizontal</property>
                    <property name="halign">end</property>
                    <property name="spacing">6</property>
                    <child>
                  <object class="GtkButton" id="copy_button">
                    <property name="icon-name">edit-copy-symbolic</property>
                    <property name="tooltip-text">Copy all results</property>
                    <style>
                      <class name="flat"/>
                    </style>
                  </object>
                </child>
                <child>
                  <object class="GtkButton" id="download_button">
                    <property name="icon-name">document-save-symbolic</property>
                    <property name="tooltip-text">Download as text file</property>
                    <style>
                      <class name="flat"/>
                    </style>
                  </object>
                </child>
              </object>
            </child>
          </object>
        </child>
        <child>
          <object class="GtkScrolledWindow" id="results_scroll">
            <property name="hscrollbar-policy">automatic</property>
            <property name="vscrollbar-policy">automatic</property>
            <property name="vexpand">true</property>
            <child>
              <object class="GtkListBox" id="results_list">
                <property name="selection-mode">none</property>
                <style>
                  <class name="boxed-list"/>
                </style>
                  </object>
                </child>
              </object>
            </child>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`

const helpUI = `<?xml version="1.0" encoding="UTF-8"?>
<interface>
  <object class="GtkDialog" id="help_dialog">
    <property name="title">Help - Shodan Facets</property>
    <property name="modal">true</property>
    <property name="default-width">600</property>
    <property name="default-height">500</property>
    <child internal-child="content_area">
      <object class="GtkBox">
        <property name="orientation">vertical</property>
        <property name="spacing">12</property>
        <property name="margin-top">16</property>
        <property name="margin-bottom">16</property>
        <property name="margin-start">16</property>
        <property name="margin-end">16</property>
        
        <child>
          <object class="GtkSearchEntry" id="facet_search_entry">
            <property name="placeholder-text">Search facets...</property>
            <property name="hexpand">true</property>
          </object>
        </child>
        
        <child>
          <object class="GtkScrolledWindow">
            <property name="hscrollbar-policy">automatic</property>
            <property name="vscrollbar-policy">automatic</property>
            <property name="vexpand">true</property>
            <child>
              <object class="GtkFlowBox" id="facets_flowbox">
                <property name="selection-mode">none</property>
                <property name="homogeneous">true</property>
                <property name="column-spacing">12</property>
                <property name="row-spacing">12</property>
                <property name="margin-top">8</property>
                <property name="margin-bottom">8</property>
                <property name="margin-start">8</property>
                <property name="margin-end">8</property>
              </object>
            </child>
          </object>
        </child>
        
      </object>
    </child>
    <child internal-child="action_area">
      <object class="GtkBox">
        <child>
          <object class="GtkButton" id="help_close_button">
            <property name="label">Close</property>
            <style>
              <class name="suggested-action"/>
            </style>
          </object>
        </child>
      </object>
    </child>
  </object>
</interface>
`

const facetsJSON = `{
  "facets": [
    {
      "name": "asn",
      "description": "Find devices by Autonomous System Number for network routing (e.g., \"AS15169\" for Google)"
    },
    {
      "name": "bitcoin.ip",
      "description": "Search for Bitcoin nodes by IP address"
    },
    {
      "name": "bitcoin.ip_count",
      "description": "Filter by the number of Bitcoin IP addresses"
    },
    {
      "name": "bitcoin.port",
      "description": "Find Bitcoin nodes running on specific ports"
    },
    {
      "name": "bitcoin.user_agent",
      "description": "Search Bitcoin nodes by their user agent string"
    },
    {
      "name": "bitcoin.version",
      "description": "Find Bitcoin nodes running specific protocol versions"
    },
    {
      "name": "city",
      "description": "Find devices in specific cities (e.g., \"New York\", \"London\", \"Tokyo\")"
    },
    {
      "name": "cloud.provider",
      "description": "Search for devices hosted on specific cloud providers (e.g., \"AWS\", \"Azure\", \"GCP\")"
    },
    {
      "name": "cloud.region",
      "description": "Find devices in specific cloud regions (e.g., \"us-east-1\", \"eu-west-1\")"
    },
    {
      "name": "cloud.service",
      "description": "Search for specific cloud services (e.g., \"EC2\", \"Compute Engine\")"
    },
    {
      "name": "country",
      "description": "Search for devices located in specific countries using country codes (e.g., US, CN, DE, RU)"
    },
    {
      "name": "cpe",
      "description": "Find devices by Common Platform Enumeration identifiers"
    },
    {
      "name": "device",
      "description": "Search for specific device types (e.g., \"router\", \"webcam\", \"printer\")"
    },
    {
      "name": "domain",
      "description": "Find devices associated with specific domains"
    },
    {
      "name": "has_screenshot",
      "description": "Filter results that have screenshots available (true/false)"
    },
    {
      "name": "hash",
      "description": "Search by various hash values (favicon, certificate, etc.)"
    },
    {
      "name": "http.component",
      "description": "Find web applications using specific components (e.g., \"jQuery\", \"Bootstrap\")"
    },
    {
      "name": "http.component_category",
      "description": "Search by web component categories (e.g., \"JavaScript Library\", \"Web Framework\")"
    },
    {
      "name": "http.dom_hash",
      "description": "Find websites with specific DOM structure hashes"
    },
    {
      "name": "http.favicon.hash",
      "description": "Search websites by their favicon hash values"
    },
    {
      "name": "http.headers_hash",
      "description": "Find websites with specific HTTP headers hash"
    },
    {
      "name": "http.html_hash",
      "description": "Search by HTML content hash values"
    },
    {
      "name": "http.robots_hash",
      "description": "Find websites by their robots.txt file hash"
    },
    {
      "name": "http.server_hash",
      "description": "Search by HTTP server response hash"
    },
    {
      "name": "http.status",
      "description": "Filter by HTTP status codes (e.g., 200, 404, 500)"
    },
    {
      "name": "http.title",
      "description": "Search by HTTP title tags from web pages (e.g., \"Welcome\", \"Login\")"
    },
    {
      "name": "http.title_hash",
      "description": "Find websites by their title hash values"
    },
    {
      "name": "http.waf",
      "description": "Search for websites protected by Web Application Firewalls"
    },
    {
      "name": "ip",
      "description": "Search for specific IP addresses or IP ranges"
    },
    {
      "name": "isp",
      "description": "Find devices by Internet Service Provider (e.g., \"Comcast\", \"Verizon\")"
    },
    {
      "name": "link",
      "description": "Search for devices with specific links or references"
    },
    {
      "name": "mongodb.database.name",
      "description": "Find MongoDB instances with specific database names"
    },
    {
      "name": "ntp.ip",
      "description": "Search NTP servers by IP address"
    },
    {
      "name": "ntp.ip_count",
      "description": "Filter by the number of NTP IP addresses"
    },
    {
      "name": "ntp.more",
      "description": "Find NTP servers with additional information available"
    },
    {
      "name": "ntp.port",
      "description": "Search NTP servers running on specific ports"
    },
    {
      "name": "org",
      "description": "Search by organization name or ISP (e.g., \"Google\", \"Amazon\")"
    },
    {
      "name": "os",
      "description": "Find devices running specific operating systems (e.g., \"Windows\", \"Linux\")"
    },
    {
      "name": "port",
      "description": "Find devices with specific open ports (e.g., 22, 80, 443)"
    },
    {
      "name": "postal",
      "description": "Search by postal/ZIP codes for geographic targeting"
    },
    {
      "name": "product",
      "description": "Search for specific software products or services (e.g., Apache, nginx)"
    },
    {
      "name": "redis.key",
      "description": "Find Redis instances with specific key patterns"
    },
    {
      "name": "region",
      "description": "Search by geographic regions or states"
    },
    {
      "name": "rsync.module",
      "description": "Find Rsync servers with specific module names"
    },
    {
      "name": "screenshot.hash",
      "description": "Search by screenshot hash values"
    },
    {
      "name": "screenshot.label",
      "description": "Find devices with specific screenshot labels or tags"
    },
    {
      "name": "snmp.contact",
      "description": "Search SNMP devices by contact information"
    },
    {
      "name": "snmp.location",
      "description": "Find SNMP devices by their configured location"
    },
    {
      "name": "snmp.name",
      "description": "Search SNMP devices by system name"
    },
    {
      "name": "ssh.cipher",
      "description": "Find SSH servers using specific cipher algorithms"
    },
    {
      "name": "ssh.fingerprint",
      "description": "Search SSH servers by their key fingerprints"
    },
    {
      "name": "ssh.hassh",
      "description": "Find SSH clients/servers by HASSH fingerprints"
    },
    {
      "name": "ssh.mac",
      "description": "Search SSH servers by MAC algorithm support"
    },
    {
      "name": "ssh.type",
      "description": "Find SSH servers by key type (e.g., RSA, ECDSA)"
    },
    {
      "name": "ssl.alpn",
      "description": "Search SSL/TLS servers by ALPN protocol support"
    },
    {
      "name": "ssl.cert.alg",
      "description": "Find SSL certificates using specific algorithms"
    },
    {
      "name": "ssl.cert.expired",
      "description": "Filter by SSL certificate expiration status (true/false)"
    },
    {
      "name": "ssl.cert.extension",
      "description": "Search SSL certificates by extension types"
    },
    {
      "name": "ssl.cert.fingerprint",
      "description": "Find SSL certificates by their fingerprints"
    },
    {
      "name": "ssl.cert.issuer.cn",
      "description": "Search SSL certificates by issuer common name"
    },
    {
      "name": "ssl.cert.pubkey.bits",
      "description": "Filter SSL certificates by public key bit length"
    },
    {
      "name": "ssl.cert.pubkey.type",
      "description": "Find SSL certificates by public key type (RSA, ECDSA)"
    },
    {
      "name": "ssl.cert.serial",
      "description": "Search SSL certificates by serial number"
    },
    {
      "name": "ssl.cert.subject.cn",
      "description": "Find SSL certificates by subject common name"
    },
    {
      "name": "ssl.chain_count",
      "description": "Filter by SSL certificate chain length"
    },
    {
      "name": "ssl.cipher.bits",
      "description": "Search SSL connections by cipher bit strength"
    },
    {
      "name": "ssl.cipher.name",
      "description": "Find SSL connections using specific cipher names"
    },
    {
      "name": "ssl.cipher.version",
      "description": "Search SSL connections by cipher version"
    },
    {
      "name": "ssl.ja3s",
      "description": "Find SSL servers by JA3S fingerprints"
    },
    {
      "name": "ssl.jarm",
      "description": "Search SSL servers by JARM fingerprints"
    },
    {
      "name": "ssl.version",
      "description": "Find SSL/TLS servers by protocol version (e.g., TLSv1.2)"
    },
    {
      "name": "state",
      "description": "Search by state or province names"
    },
    {
      "name": "tag",
      "description": "Find devices with specific Shodan tags"
    },
    {
      "name": "telnet.do",
      "description": "Search Telnet servers by DO option negotiations"
    },
    {
      "name": "telnet.dont",
      "description": "Find Telnet servers by DONT option negotiations"
    },
    {
      "name": "telnet.option",
      "description": "Search Telnet servers by supported options"
    },
    {
      "name": "telnet.will",
      "description": "Find Telnet servers by WILL option negotiations"
    },
    {
      "name": "telnet.wont",
      "description": "Search Telnet servers by WONT option negotiations"
    },
    {
      "name": "uptime",
      "description": "Filter devices by their reported uptime"
    },
    {
      "name": "version",
      "description": "Search by software or protocol version numbers"
    },
    {
      "name": "vuln",
      "description": "Find devices with specific vulnerabilities (CVE IDs)"
    },
    {
      "name": "vuln.verified",
      "description": "Filter by verified vulnerability status (true/false)"
    }
  ]
}`

const aboutUI = `<?xml version="1.0" encoding="UTF-8"?>
<interface>
  <object class="GtkAboutDialog" id="about_dialog">
    <property name="program-name">shef</property>
    <property name="version">1.0.0</property>
    <property name="comments">bring shodan facets into your terminal without API key.</property>
    <property name="website">https://github.com/1hehaq/shef</property>
    <property name="website-label">GitHub Repository</property>
    <property name="authors">1hehaq</property>
    <property name="copyright">© 2025 haq</property>
    <property name="license-type">mit</property>
    <property name="logo-icon-name">website</property>
  </object>
</interface>
`
