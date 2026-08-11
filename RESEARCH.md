# VM Cloud-Init Webapp — Language Selection Research

## Verdict: Go is the Clear Winner

---

## Comparison Overview

| Criteria | Go | Java | PHP |
|---|---|---|---|
| **Stars (main binding)** | 1.1k | 58 | 39 |
| **Maintained** | Actively (2026 commits) | Stale mirror | PECL: no releases |
| **Native C required** | No (pure Go RPC) | Yes (JNA/FIFO) | Yes (compiled ext) |
| **SSH transport** | Built-in | Via libvirt C | Via libvirt C |
| **ISO creation** | `go-diskfs` (689 ⭐) | None | None |
| **Ecosystem proof** | kubevirt (7k ⭐), terraform-provider-libvirt (1.8k ⭐) | Apache CloudStack | — |
| **Deployment** | Single static binary | JVM + dependencies | PHP-FPM + extensions |
| **Container-friendly** | Excellent | Heavy (JVM overhead) | Moderate |

---

## Detailed Analysis

### Go

**Libvirt binding: `github.com/digitalocean/go-libvirt`**
- Pure Go implementation — speaks libvirt's native RPC/XDR protocol directly
- Zero C dependencies: no `libvirt-dev` needed to compile
- Full API coverage via code generators from libvirt source
- Supports all connection URIs: local, SSH, TLS, TCP, Unix socket
- Actively maintained by DigitalOcean, 352 commits
- Apache 2.0 license

**Ecosystem for cloud-init workflow:**
- `go-diskfs` (689 ⭐) — native ISO9660/ filesystem creation for NoCloud cloud-init ISOs
- `quiso` (27 ⭐) — dedicated cloud-init ISO builder CLI
- `sigs.k8s.io/yaml` — YAML templating for user-data/meta-data generation
- Standard library `net/http` for downloading base images

**Real-world proof:**
- `kubevirt` (7k ⭐) — Kubernetes VM runtime
- `terraform-provider-libvirt` (1.8k ⭐) — Terraform provider for KVM
- `quickpassthrough` (924 ⭐) — GPU passthrough tooling

**Example connection:**
```go
uri, _ := url.Parse("qemu+ssh://user@host/system")
conn, err := libvirt.ConnectToURI(uri)
```

### Java

**Libvirt binding: `libvirt/libvirt-java`**
- Official binding, uses JNA (Java Native Access) over FIFO socket
- Requires native libvirt library installed on the system
- Read-only GitHub mirror; issues go to GitLab
- 58 stars, lower community engagement for infrastructure tooling
- No ecosystem support for disk image creation or cloud-init

**Drawbacks:**
- JVM startup overhead and memory footprint
- No native ISO/disk image manipulation libraries
- Heavier deployment model

### PHP

**Libvirt binding: PECL `libvirt` extension**
- Must be compiled against libvirt C headers
- PECL page shows **no releases available**
- Official repo is a read-only mirror (39 stars, last update 2026-06-09 but content is stale)
- Most PHP libvirt projects are archived (last updates 2010-2019)

**Drawbacks:**
- Requires compiling native extensions on each target system
- Wrong paradigm for infrastructure tooling
- No disk image or cloud-init libraries
- Request/response model not suited for long-running VM operations

---

## Recommended Architecture (Go)

```
┌─────────────────────────────────────────────────────────────┐
│                     Web Application                         │
│  (Gin / Echo / Fiber for HTTP API + frontend framework)     │
├──────────────┬──────────────┬───────────────────────────────┤
│  Libvirt     │  Image       │  Cloud-Init Config            │
│  Service     │  Manager     │  Generator                    │
│              │              │                               │
│  • Connect   │  • Download  │  • user-data (YAML)           │
│  • SSH       │  • Base      │  • meta-data                  │
│  • TLS       │    images    │  • ISO creation (go-diskfs)   │
│  • List VMs  │  • Copy to   │  • Upload to libvirt storage  │
│  • Create    │    storage   │                               │
│  • Delete    │              │                               │
├──────────────┴──────────────┴───────────────────────────────┤
│                     libvirt (qemu:///system)                │
└─────────────────────────────────────────────────────────────┘
```

## Key Dependencies

| Purpose | Package | Stars |
|---|---|---|
| Libvirt RPC | `github.com/digitalocean/go-libvirt` | 1.1k |
| ISO creation | `github.com/diskfs/go-diskfs` | 689 |
| YAML parsing | `gopkg.in/yaml.v3` or `sigs.k8s.io/yaml` | — |
| HTTP client | stdlib `net/http` | — |
| Web framework | `github.com/gofiber/fiber` or `github.com/gin-gonic/gin` | 30k+ |

## Conclusion

Go dominates for this use case across every criterion:
1. **Best libvirt binding** — pure Go, no C deps, actively maintained
2. **Complete ecosystem** — ISO creation, YAML templating, HTTP clients all native
3. **Proven in production** — kubevirt, DigitalOcean, Terraform provider
4. **Simple deployment** — single binary, no runtime dependencies
5. **SSH support** — built-in via dialers

Java and PHP both require native libvirt C bindings, have stale or minimal community engagement, and lack supporting libraries for the cloud-init workflow.
