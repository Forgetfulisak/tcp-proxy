Proof of concept for creating a tcp-proxy for www.nrk.no

## Usage:

- 1
    Add the following to  `C:\Windows\System32\drivers\etc\hosts` (windows) or `/etc/hosts` (linux).

    ```
    127.0.0.1 nrk.no
    127.0.0.1 www.nrk.no
    ```
    This forces chrome to connect to our proxy instead of the real nrk-servers. You may have to flush dns and/or restart the browser for the changes to take effect.


- 2
  Build/run the proxy with 
  `go run .`

- 3 Open a browser and go to www.nrk.no.
    The browser will (rightly) warn you that the connection is unsafe. However, if you ignore the warning, you will see the proxy print out all requests sent by the browser
