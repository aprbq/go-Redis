# go-redis

โปรเจกต์ทดลอง Go Fiber server พร้อม stack สำหรับ load testing และ monitoring — ตั้งใจใช้ศึกษาผลของการทำ caching ด้วย Redis ต่อ performance ของ API

## ภาพรวม

```
k6 (load test) ──▶ Go Fiber server (:8000) ──▶ Redis (:6379)
      │
      ▼
  InfluxDB (:8086) ──▶ Grafana (:3000)
```

- **Go Fiber server** ([main.go](main.go)) — HTTP server ด้วย [Fiber v3](https://gofiber.io) มี endpoint `GET /hello` ที่จำลอง latency 10ms (ยังไม่ได้เชื่อมต่อ Redis — เตรียมไว้สำหรับทดลอง caching)
- **Redis** — รันผ่าน Docker ใช้ config จาก [config/redis.conf](config/redis.conf) (เปิด AOF persistence, ปิด RDB snapshot) เก็บข้อมูลไว้ที่ `data/redis`
- **k6** — load testing script อยู่ที่ [scripts/test.js](scripts/test.js) (5 VUs, 5 วินาที ยิงไปที่ `/hello`) ส่งผลลัพธ์เข้า InfluxDB
- **InfluxDB 1.8** — เก็บ metrics จาก k6 (database ชื่อ `k6`)
- **Grafana** — visualize ผล load test ที่ http://localhost:3000 (เปิด anonymous access เป็น Admin)

## วิธีใช้งาน

### 1. Start stack (Redis, InfluxDB, Grafana)

```bash
docker compose up -d redis influxdb grafana
```

### 2. Run Go server

```bash
go run main.go
```

Server จะฟังที่ http://localhost:8000 — ทดสอบด้วย:

```bash
curl http://localhost:8000/hello
```

### 3. Run load test ด้วย k6

```bash
docker compose run --rm k6 run /scripts/test.js
```

k6 จะยิง request ไปที่ server บนเครื่อง host (ผ่าน `host.docker.internal:8000`) และส่ง metrics เข้า InfluxDB จากนั้นดูผลได้ใน Grafana โดยเพิ่ม InfluxDB data source (`http://influxdb:8086`, database `k6`) แล้ว import dashboard สำหรับ k6

## โครงสร้างโปรเจกต์

```
├── main.go              # Fiber server
├── docker-compose.yml   # Redis, k6, InfluxDB, Grafana
├── config/
│   └── redis.conf       # Redis config (AOF enabled)
├── scripts/
│   └── test.js          # k6 load test script
└── data/                # Volume data ของ Redis / InfluxDB / Grafana
```

## Requirements

- Go 1.25+
- Docker + Docker Compose
