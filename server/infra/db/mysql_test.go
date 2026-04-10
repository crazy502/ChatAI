package db

import (
	"testing"
	"time"

	"server/infra/config"
)

func TestNormalizePoolConfigAppliesDefaults(t *testing.T) {
	pool := normalizePoolConfig(config.MysqlPoolConfig{})

	if pool.MaxIdleConns != defaultMaxIdleConns {
		t.Fatalf("expected default max idle conns %d, got %d", defaultMaxIdleConns, pool.MaxIdleConns)
	}
	if pool.MaxOpenConns != defaultMaxOpenConns {
		t.Fatalf("expected default max open conns %d, got %d", defaultMaxOpenConns, pool.MaxOpenConns)
	}
	if pool.ConnMaxLifetime != defaultConnMaxLifetime {
		t.Fatalf("expected default conn max lifetime %s, got %s", defaultConnMaxLifetime, pool.ConnMaxLifetime)
	}
	if pool.ConnMaxIdleTime != defaultConnMaxIdleTime {
		t.Fatalf("expected default conn max idle time %s, got %s", defaultConnMaxIdleTime, pool.ConnMaxIdleTime)
	}
}

func TestNormalizePoolConfigUsesConfiguredValues(t *testing.T) {
	pool := normalizePoolConfig(config.MysqlPoolConfig{
		MaxIdleConns:           7,
		MaxOpenConns:           42,
		ConnMaxLifetimeMinutes: 90,
		ConnMaxIdleTimeMinutes: 15,
	})

	if pool.MaxIdleConns != 7 {
		t.Fatalf("expected configured max idle conns 7, got %d", pool.MaxIdleConns)
	}
	if pool.MaxOpenConns != 42 {
		t.Fatalf("expected configured max open conns 42, got %d", pool.MaxOpenConns)
	}
	if pool.ConnMaxLifetime != 90*time.Minute {
		t.Fatalf("expected configured conn max lifetime %s, got %s", 90*time.Minute, pool.ConnMaxLifetime)
	}
	if pool.ConnMaxIdleTime != 15*time.Minute {
		t.Fatalf("expected configured conn max idle time %s, got %s", 15*time.Minute, pool.ConnMaxIdleTime)
	}
}

func TestBuildReplicaConnectionConfigFallsBackToPrimary(t *testing.T) {
	cfg := &config.Config{
		MysqlConfig: config.MysqlConfig{
			MysqlHost:         "primary-host",
			MysqlPort:         3306,
			MysqlUser:         "writer",
			MysqlPassword:     "writer-secret",
			MysqlDatabaseName: "agentgo",
			MysqlCharset:      "utf8mb4",
			Pool: config.MysqlPoolConfig{
				MaxIdleConns:           10,
				MaxOpenConns:           100,
				ConnMaxLifetimeMinutes: 60,
				ConnMaxIdleTimeMinutes: 10,
			},
			Replica: config.MysqlReplicaConfig{
				Enabled: true,
				Host:    "replica-host",
				Pool: config.MysqlPoolConfig{
					MaxIdleConns: 5,
				},
			},
		},
	}

	replicaCfg, ok := buildReplicaConnectionConfig(cfg)
	if !ok {
		t.Fatal("expected replica config to be enabled")
	}

	if replicaCfg.Host != "replica-host" {
		t.Fatalf("expected replica host replica-host, got %s", replicaCfg.Host)
	}
	if replicaCfg.Port != 3306 {
		t.Fatalf("expected replica port to fall back to primary 3306, got %d", replicaCfg.Port)
	}
	if replicaCfg.User != "writer" {
		t.Fatalf("expected replica user to fall back to primary writer, got %s", replicaCfg.User)
	}
	if replicaCfg.Password != "writer-secret" {
		t.Fatalf("expected replica password to fall back to primary, got %s", replicaCfg.Password)
	}
	if replicaCfg.DatabaseName != "agentgo" {
		t.Fatalf("expected replica database to fall back to primary agentgo, got %s", replicaCfg.DatabaseName)
	}
	if replicaCfg.Charset != "utf8mb4" {
		t.Fatalf("expected replica charset to fall back to utf8mb4, got %s", replicaCfg.Charset)
	}
	if replicaCfg.Pool.MaxIdleConns != 5 {
		t.Fatalf("expected replica pool override max idle conns 5, got %d", replicaCfg.Pool.MaxIdleConns)
	}
	if replicaCfg.Pool.MaxOpenConns != 100 {
		t.Fatalf("expected replica pool max open conns to inherit 100, got %d", replicaCfg.Pool.MaxOpenConns)
	}
}

func TestBuildReplicaConnectionConfigDisabled(t *testing.T) {
	cfg := &config.Config{
		MysqlConfig: config.MysqlConfig{
			MysqlHost:         "primary-host",
			MysqlPort:         3306,
			MysqlUser:         "writer",
			MysqlPassword:     "writer-secret",
			MysqlDatabaseName: "agentgo",
			MysqlCharset:      "utf8mb4",
		},
	}

	if _, ok := buildReplicaConnectionConfig(cfg); ok {
		t.Fatal("expected replica config to be disabled")
	}
}
