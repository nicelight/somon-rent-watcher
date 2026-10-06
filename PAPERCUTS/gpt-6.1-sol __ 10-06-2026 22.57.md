# Papercuts

- Docker daemon rejected a nested read-only file bind over the project directory bind: mountpoint outside rootfs, exit 125. This setup failure is not RED. The identical baseline probe was copied into internal/filter for the network-disabled package test, then removed before final gates; no production code affected.
