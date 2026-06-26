package cli

import (
	"context"
	"fmt"
	"strings"

	"zabbix-maint/internal/config"
	"zabbix-maint/internal/version"
	"zabbix-maint/pkg/zabbix"
)

// OneShotHandler 一次性命令处理器
type OneShotHandler struct {
	client zabbix.ZabbixOperator
}

// NewOneShotHandlerWithInstance 创建连接实例的一次性命令处理器
func NewOneShotHandlerWithInstance(ctx context.Context, instanceName string) (*OneShotHandler, error) {
	cfg, err := config.Load(instanceName)
	if err != nil {
		return nil, err
	}

	router := version.NewRouter()
	client, _, err := router.Connect(ctx, cfg)
	if err != nil {
		return nil, err
	}

	return &OneShotHandler{client: client}, nil
}

// Execute 执行一次性命令
func (h *OneShotHandler) Execute(ctx context.Context, args []string) error {
	if h.client == nil {
		return fmt.Errorf("client not initialized")
	}
	if len(args) < 2 {
		return fmt.Errorf("usage: <resource> <action> [flags...]")
	}

	resource := args[0]
	action := args[1]
	flags := parseFlags(args[2:])

	switch resource {
	case "user":
		return h.handleUser(ctx, action, flags)
	case "usergroup":
		return h.handleUserGroup(ctx, action, flags)
	case "host":
		return h.handleHost(ctx, action, flags)
	case "hostgroup":
		return h.handleHostGroup(ctx, action, flags)
	case "dashboard":
		return h.handleDashboard(ctx, action, flags)
	case "system":
		return h.handleSystem(ctx, action, flags)
	default:
		return fmt.Errorf("unknown resource: %s", resource)
	}
}

func parseFlags(args []string) map[string]string {
	flags := make(map[string]string)
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			key := args[i][2:]
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				flags[key] = args[i+1]
				i++
			} else {
				flags[key] = "true"
			}
		}
	}
	return flags
}

func (h *OneShotHandler) handleUser(ctx context.Context, action string, flags map[string]string) error {
	switch action {
	case "list":
		users, err := h.client.UserList(ctx, "")
		if err != nil {
			return err
		}
		fmt.Println("ID\tALIAS\tNAME\tSTATUS")
		for _, u := range users {
			status := "Enabled"
			if !u.Enabled {
				status = "Disabled"
			}
			fmt.Printf("%s\t%s\t%s\t%s\n", u.ID, u.Alias, u.Name, status)
		}
		return nil
	case "create":
		req := zabbix.UserCreateReq{
			Alias:  flags["alias"],
			Name:   flags["name"],
			Passwd: flags["password"],
		}
		if gid, ok := flags["group"]; ok {
			req.UserGroupIDs = []string{gid}
		}
		if rid, ok := flags["role"]; ok {
			req.RoleID = rid
		}
		id, err := h.client.UserCreate(ctx, req)
		if err != nil {
			return err
		}
		fmt.Printf("User created: %s\n", id)
		return nil
	case "disable":
		return h.client.UserDisable(ctx, flags["id"])
	case "enable":
		return h.client.UserEnable(ctx, flags["id"])
	case "passwd":
		return h.client.UserUpdatePassword(ctx, flags["id"], flags["password"])
	case "delete":
		return h.client.UserDelete(ctx, flags["id"])
	case "detail":
		u, err := h.client.UserDetail(ctx, flags["id"])
		if err != nil {
			return err
		}
		fmt.Printf("ID: %s\nAlias: %s\nName: %s\nEnabled: %v\n", u.ID, u.Alias, u.Name, u.Enabled)
		return nil
	default:
		return fmt.Errorf("unknown user action: %s", action)
	}
}

func (h *OneShotHandler) handleUserGroup(ctx context.Context, action string, flags map[string]string) error {
	switch action {
	case "list":
		groups, err := h.client.UserGroupList(ctx, "")
		if err != nil {
			return err
		}
		fmt.Println("ID\tNAME")
		for _, g := range groups {
			fmt.Printf("%s\t%s\n", g.ID, g.Name)
		}
		return nil
	case "create":
		req := zabbix.UserGroupCreateReq{
			Name: flags["name"],
		}
		id, err := h.client.UserGroupCreate(ctx, req)
		if err != nil {
			return err
		}
		fmt.Printf("UserGroup created: %s\n", id)
		return nil
	case "disable":
		return h.client.UserGroupDisable(ctx, flags["id"])
	case "enable":
		return h.client.UserGroupEnable(ctx, flags["id"])
	case "clear":
		return h.client.UserGroupClear(ctx, flags["id"])
	case "delete":
		return h.client.UserGroupDelete(ctx, flags["id"])
	default:
		return fmt.Errorf("unknown usergroup action: %s", action)
	}
}

func (h *OneShotHandler) handleHost(ctx context.Context, action string, flags map[string]string) error {
	switch action {
	case "list":
		hosts, err := h.client.HostList(ctx, "")
		if err != nil {
			return err
		}
		fmt.Println("ID\tHOST\tNAME")
		for _, h := range hosts {
			fmt.Printf("%s\t%s\t%s\n", h.ID, h.Host, h.Name)
		}
		return nil
	case "clone":
		req := zabbix.HostCloneReq{
			SourceHostID: flags["src"],
			NewHostName:  flags["name"],
			NewHostIP:    flags["ip"],
		}
		id, err := h.client.HostFullClone(ctx, req)
		if err != nil {
			return err
		}
		fmt.Printf("Host cloned: %s\n", id)
		return nil
	case "detail":
		host, err := h.client.HostDetail(ctx, flags["id"])
		if err != nil {
			return err
		}
		fmt.Printf("ID: %s\nHost: %s\nName: %s\n", host.ID, host.Host, host.Name)
		return nil
	default:
		return fmt.Errorf("unknown host action: %s", action)
	}
}

func (h *OneShotHandler) handleHostGroup(ctx context.Context, action string, flags map[string]string) error {
	switch action {
	case "list":
		groups, err := h.client.HostGroupList(ctx, "")
		if err != nil {
			return err
		}
		fmt.Println("ID\tNAME")
		for _, g := range groups {
			fmt.Printf("%s\t%s\n", g.ID, g.Name)
		}
		return nil
	case "create":
		id, err := h.client.HostGroupCreate(ctx, flags["name"])
		if err != nil {
			return err
		}
		fmt.Printf("HostGroup created: %s\n", id)
		return nil
	case "clear":
		return h.client.HostGroupClear(ctx, flags["id"])
	case "delete":
		return h.client.HostGroupDelete(ctx, flags["id"])
	default:
		return fmt.Errorf("unknown hostgroup action: %s", action)
	}
}

func (h *OneShotHandler) handleDashboard(ctx context.Context, action string, flags map[string]string) error {
	switch action {
	case "list":
		dashs, err := h.client.DashboardList(ctx, "")
		if err != nil {
			return err
		}
		fmt.Println("ID\tNAME\tTYPE")
		for _, d := range dashs {
			fmt.Printf("%s\t%s\t%s\n", d.ID, d.Name, d.Type)
		}
		return nil
	case "create":
		req := zabbix.DashboardCreateReq{
			Name: flags["name"],
		}
		id, err := h.client.DashboardCreate(ctx, req)
		if err != nil {
			return err
		}
		fmt.Printf("Dashboard created: %s\n", id)
		return nil
	default:
		return fmt.Errorf("unknown dashboard action: %s", action)
	}
}

func (h *OneShotHandler) handleSystem(ctx context.Context, action string, flags map[string]string) error {
	switch action {
	case "version":
		ver, err := h.client.APIVersion(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("API Version: %s\n", ver)
		return nil
	case "status":
		status, err := h.client.ServerStatus(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("Status: %+v\n", status)
		return nil
	default:
		return fmt.Errorf("unknown system action: %s", action)
	}
}
