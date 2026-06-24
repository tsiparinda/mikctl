CREATE VIEW routers_by_group AS 			
SELECT r.*
				FROM routers r
                                    join devices d on r.device_id=d.id
                                    join router_groups rg on rg.device_id=d.id
                                    join groups g on g.id=rg.group_id
				WHERE r.active = 1
				and g.name = 'v7'
				ORDER BY r.name;
CREATE TABLE passwords (id INTEGER PRIMARY KEY, password TEXT NOT NULL);
CREATE TABLE devices (id INTEGER PRIMARY KEY, ros_version TEXT, model TEXT, serial_number TEXT, created_at TEXT DEFAULT (date('now')) NOT NULL, comment TEXT);
CREATE TABLE script_runs (id INTEGER PRIMARY KEY, router_id INTEGER NOT NULL, tag TEXT, status TEXT NOT NULL, started_at TEXT NOT NULL, finished_at TEXT, script_name TEXT NOT NULL, message TEXT, FOREIGN KEY (router_id) REFERENCES routers (id));
CREATE UNIQUE INDEX idx_devices_serial ON devices (serial_number);
CREATE INDEX idx_scripts_run_tag ON script_runs (tag);
CREATE INDEX idx_scripts_run_status ON script_runs (status);
CREATE INDEX idx_scripts_run_started_at ON script_runs (started_at);
CREATE TABLE router_groups (router_id INTEGER NOT NULL, name NOT NULL, PRIMARY KEY (router_id), FOREIGN KEY (router_id) REFERENCES routers (id));
CREATE UNIQUE INDEX idx_router_groups_uniq ON router_groups (router_id, name);
CREATE VIEW routers_full AS SELECT 
r.id,
r.name,
coalesce(a.ip, '') ip,
coalesce(r.ssh_user, '') ssh_user,
coalesce(r.site, '') site,
COALESCE(r.device_id,0) as device_id,
r.fMain,
d.ros_version,
d.model,
coalesce(d.serial_number, '') serial_number
FROM routers r 
join routers_addresses a on r.id=a.router_id and a.fMain=1 and r.active = 1
left join devices d on d.id=r.device_id
/* routers_full(id,name,ip,ssh_user,site,device_id,fMain,ros_version,model,serial_number) */;
CREATE VIEW routers_and_device AS select r.*, d.model, d.ros_version, a.ip, a.fMain
from routers r
left join devices d on r.device_id=d.id
left join routers_addresses a on a.router_id=r.id
order by r.name
/* routers_and_device(id,device_id,name,ssh_user,site,installed_at,active,password_id,fMain,parent_router_id,last_seen_at,comment,prescript_id,model,ros_version,ip,"fMain:1") */;
CREATE TABLE routers_addresses (router_id INTEGER, ip TEXT NOT NULL, comment TEXT, fMain INTEGER NOT NULL DEFAULT (0) CHECK (fMain IN (0, 1)), PRIMARY KEY (router_id, ip), FOREIGN KEY (router_id) REFERENCES routers (id) ON DELETE CASCADE);
CREATE UNIQUE INDEX ix_routers_addresses_ip ON routers_addresses (ip);
CREATE UNIQUE INDEX ix_routers_addresses_main ON routers_addresses (router_id) WHERE fMain = 1;
CREATE TABLE routers (id INTEGER PRIMARY KEY, device_id INTEGER, name TEXT NOT NULL, ssh_user TEXT NOT NULL, site TEXT, installed_at TEXT DEFAULT (date('now')) NOT NULL, active INTEGER DEFAULT 1, password_id INTEGER NOT NULL, fMain INTEGER NOT NULL DEFAULT 0, parent_router_id INTEGER REFERENCES routers (id), last_seen_at TEXT, comment TEXT, prescript_id INTEGER DEFAULT (1) NOT NULL, FOREIGN KEY (password_id) REFERENCES passwords (id), FOREIGN KEY (prescript_id) REFERENCES prescripts (id));
CREATE INDEX idx_routers_parent ON routers (parent_router_id);
CREATE UNIQUE INDEX idx_routers_device_id ON routers (device_id);
CREATE UNIQUE INDEX idx_routers_name ON routers (name);
CREATE INDEX idx_routers_last_seen ON routers (last_seen_at);
CREATE INDEX idx_routers_site ON routers (site);
CREATE INDEX idx_routers_main ON routers (fMain);
CREATE TABLE prescripts (id INTEGER PRIMARY KEY, prescript TEXT NOT NULL);
