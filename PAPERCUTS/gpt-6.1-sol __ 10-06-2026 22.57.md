# Papercuts

- Docker daemon rejected a nested read-only file bind over the project directory bind: mountpoint outside rootfs, exit 125. This setup failure is not RED. The identical baseline probe was copied into internal/filter for the network-disabled package test, then removed before final gates; no production code affected.

- HEAD changed externally during execution from 7782bc1 to d99e3a7, tracking prepared protocol and the temporary baseline package test. The intended cleanup therefore appears as a tracked deletion; original probe remains under .tasks. Executor made no Git mutations and offers no cached/reuse gate proof. Owner informed.
