package server

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	pb "github.com/magefree/mage-server-go/pkg/proto/mage/v1"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// HTTPJSONHandler provides simple HTTP/JSON endpoints for gRPC methods
type HTTPJSONHandler struct {
	mageServer pb.MageServerServer
	logger     *zap.Logger
}

// NewHTTPJSONHandler creates a new HTTP/JSON handler
func NewHTTPJSONHandler(mageServer pb.MageServerServer, logger *zap.Logger) *HTTPJSONHandler {
	return &HTTPJSONHandler{
		mageServer: mageServer,
		logger:     logger,
	}
}

// ServeHTTP handles HTTP/JSON requests
func (h *HTTPJSONHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log all incoming requests for debugging
	h.logger.Debug("HTTP request received",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.String("remote_addr", r.RemoteAddr),
	)

	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Grpc-Web")

	// Handle preflight
	if r.Method == http.MethodOptions {
		h.logger.Debug("CORS preflight request")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only handle POST requests for RPC
	if r.Method != http.MethodPost {
		h.logger.Warn("non-POST request rejected", zap.String("method", r.Method))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract method name from path: /mage.v1.MageServer/MethodName
	path := strings.TrimPrefix(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] != "mage.v1.MageServer" {
		h.logger.Warn("invalid path format", zap.String("path", r.URL.Path))
		http.Error(w, "Invalid path format. Expected: /mage.v1.MageServer/MethodName", http.StatusBadRequest)
		return
	}

	methodName := parts[1]
	ctx := r.Context()

	h.logger.Info("handling RPC request",
		zap.String("rpc_method", methodName),
		zap.String("remote_addr", r.RemoteAddr),
	)

	// Route to the appropriate handler
	switch methodName {
	// Authentication & Connection
	case "AuthRegister":
		h.handleAuthRegister(ctx, w, r)
	case "AuthSendTokenToEmail":
		h.handleAuthSendTokenToEmail(ctx, w, r)
	case "AuthResetPassword":
		h.handleAuthResetPassword(ctx, w, r)
	case "ConnectUser":
		h.handleConnectUser(ctx, w, r)
	case "ConnectAdmin":
		h.handleConnectAdmin(ctx, w, r)
	case "ConnectSetUserData":
		h.handleConnectSetUserData(ctx, w, r)
	case "Ping":
		h.handlePing(ctx, w, r)

	// Server Info
	case "GetServerState":
		h.handleGetServerState(ctx, w, r)
	case "ServerGetPromotionMessages":
		h.handleServerGetPromotionMessages(ctx, w, r)
	case "ServerAddFeedbackMessage":
		h.handleServerAddFeedbackMessage(ctx, w, r)

	// Room/Lobby
	case "ServerGetMainRoomId":
		h.handleServerGetMainRoomId(ctx, w, r)
	case "RoomGetUsers":
		h.handleRoomGetUsers(ctx, w, r)
	case "RoomGetFinishedMatches":
		h.handleRoomGetFinishedMatches(ctx, w, r)
	case "RoomGetAllTables":
		h.handleRoomGetAllTables(ctx, w, r)
	case "RoomGetTableById":
		h.handleRoomGetTableById(ctx, w, r)

	// Table Management
	case "RoomCreateTable":
		h.handleRoomCreateTable(ctx, w, r)
	case "RoomCreateTournament":
		h.handleRoomCreateTournament(ctx, w, r)
	case "RoomJoinTable":
		h.handleRoomJoinTable(ctx, w, r)
	case "RoomJoinTournament":
		h.handleRoomJoinTournament(ctx, w, r)
	case "RoomLeaveTableOrTournament":
		h.handleRoomLeaveTableOrTournament(ctx, w, r)
	case "RoomWatchTable":
		h.handleRoomWatchTable(ctx, w, r)
	case "RoomWatchTournament":
		h.handleRoomWatchTournament(ctx, w, r)
	case "TableSwapSeats":
		h.handleTableSwapSeats(ctx, w, r)
	case "TableRemove":
		h.handleTableRemove(ctx, w, r)
	case "TableIsOwner":
		h.handleTableIsOwner(ctx, w, r)

	// Deck Management
	case "DeckSubmit":
		h.handleDeckSubmit(ctx, w, r)
	case "DeckSave":
		h.handleDeckSave(ctx, w, r)
	case "DeckList":
		h.handleDeckList(ctx, w, r)
	case "DeckDelete":
		h.handleDeckDelete(ctx, w, r)
	case "DeckGet":
		h.handleDeckGet(ctx, w, r)

	// Game Execution
	case "GameJoin":
		h.handleGameJoin(ctx, w, r)
	case "GameWatchStart":
		h.handleGameWatchStart(ctx, w, r)
	case "GameWatchStop":
		h.handleGameWatchStop(ctx, w, r)
	case "GameGetView":
		h.handleGameGetView(ctx, w, r)
	case "SendPlayerUUID":
		h.handleSendPlayerUUID(ctx, w, r)
	case "SendPlayerString":
		h.handleSendPlayerString(ctx, w, r)
	case "SendPlayerBoolean":
		h.handleSendPlayerBoolean(ctx, w, r)
	case "SendPlayerInteger":
		h.handleSendPlayerInteger(ctx, w, r)
	case "SendPlayerManaType":
		h.handleSendPlayerManaType(ctx, w, r)
	case "SendPlayerAction":
		h.handleSendPlayerAction(ctx, w, r)
	case "MatchStart":
		h.handleMatchStart(ctx, w, r)
	case "MatchQuit":
		h.handleMatchQuit(ctx, w, r)
	case "SendSpecialAction":
		h.handleSendSpecialAction(ctx, w, r)

	// Draft
	case "DraftJoin":
		h.handleDraftJoin(ctx, w, r)
	case "SendDraftCardPick":
		h.handleSendDraftCardPick(ctx, w, r)
	case "SendDraftCardMark":
		h.handleSendDraftCardMark(ctx, w, r)
	case "DraftSetBoosterLoaded":
		h.handleDraftSetBoosterLoaded(ctx, w, r)
	case "DraftQuit":
		h.handleDraftQuit(ctx, w, r)

	// Tournament
	case "TournamentJoin":
		h.handleTournamentJoin(ctx, w, r)
	case "TournamentStart":
		h.handleTournamentStart(ctx, w, r)
	case "TournamentQuit":
		h.handleTournamentQuit(ctx, w, r)
	case "TournamentFindById":
		h.handleTournamentFindById(ctx, w, r)

	// Chat
	case "ChatJoin":
		h.handleChatJoin(ctx, w, r)
	case "ChatLeave":
		h.handleChatLeave(ctx, w, r)
	case "ChatSendMessage":
		h.handleChatSendMessage(ctx, w, r)
	case "ChatFindByTable":
		h.handleChatFindByTable(ctx, w, r)
	case "ChatFindByGame":
		h.handleChatFindByGame(ctx, w, r)
	case "ChatFindByTournament":
		h.handleChatFindByTournament(ctx, w, r)
	case "ChatFindByRoom":
		h.handleChatFindByRoom(ctx, w, r)

	// Match History
	case "GetMatchHistory":
		h.handleGetMatchHistory(ctx, w, r)
	case "GetMatchById":
		h.handleGetMatchById(ctx, w, r)

	// Replay
	case "ReplayInit":
		h.handleReplayInit(ctx, w, r)
	case "ReplayStart":
		h.handleReplayStart(ctx, w, r)
	case "ReplayStop":
		h.handleReplayStop(ctx, w, r)
	case "ReplayNext":
		h.handleReplayNext(ctx, w, r)
	case "ReplayPrevious":
		h.handleReplayPrevious(ctx, w, r)
	case "ReplaySkipForward":
		h.handleReplaySkipForward(ctx, w, r)

	// Admin
	case "AdminGetUsers":
		h.handleAdminGetUsers(ctx, w, r)
	case "AdminDisconnectUser":
		h.handleAdminDisconnectUser(ctx, w, r)
	case "AdminMuteUser":
		h.handleAdminMuteUser(ctx, w, r)
	case "AdminLockUser":
		h.handleAdminLockUser(ctx, w, r)
	case "AdminActivateUser":
		h.handleAdminActivateUser(ctx, w, r)
	case "AdminToggleActivateUser":
		h.handleAdminToggleActivateUser(ctx, w, r)
	case "AdminEndUserSession":
		h.handleAdminEndUserSession(ctx, w, r)
	case "AdminTableRemove":
		h.handleAdminTableRemove(ctx, w, r)
	case "AdminSendBroadcastMessage":
		h.handleAdminSendBroadcastMessage(ctx, w, r)

	default:
		h.logger.Warn("unsupported method", zap.String("method", methodName))
		http.Error(w, "Method not supported in HTTP/JSON mode: "+methodName, http.StatusNotImplemented)
	}
}

// handleAuthRegister handles the AuthRegister method
func (h *HTTPJSONHandler) handleAuthRegister(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AuthRegisterRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AuthRegister(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleConnectUser handles the ConnectUser method
func (h *HTTPJSONHandler) handleConnectUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ConnectUserRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ConnectUser(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handlePing handles the Ping method
func (h *HTTPJSONHandler) handlePing(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.PingRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.Ping(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleGetServerState handles the GetServerState method
func (h *HTTPJSONHandler) handleGetServerState(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.GetServerStateRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.GetServerState(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleServerGetMainRoomId handles the ServerGetMainRoomId method
func (h *HTTPJSONHandler) handleServerGetMainRoomId(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ServerGetMainRoomIdRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ServerGetMainRoomId(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomGetAllTables handles the RoomGetAllTables method
func (h *HTTPJSONHandler) handleRoomGetAllTables(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomGetAllTablesRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomGetAllTables(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDeckList handles the DeckList method
func (h *HTTPJSONHandler) handleDeckList(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DeckListRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DeckList(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDeckGet handles the DeckGet method
func (h *HTTPJSONHandler) handleDeckGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DeckGetRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DeckGet(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDeckDelete handles the DeckDelete method
func (h *HTTPJSONHandler) handleDeckDelete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DeckDeleteRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DeckDelete(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDeckSave handles the DeckSave method
func (h *HTTPJSONHandler) handleDeckSave(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DeckSaveRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DeckSave(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// camelToSnake converts camelCase to snake_case for proto field name compatibility
func camelToSnake(s string) string {
	var result []byte
	for i, c := range s {
		if c >= 'A' && c <= 'Z' {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, byte(c-'A'+'a'))
		} else {
			result = append(result, byte(c))
		}
	}
	return string(result)
}

// convertKeysToSnakeCase recursively converts all map keys from camelCase to snake_case
func convertKeysToSnakeCase(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			snakeKey := camelToSnake(k)
			result[snakeKey] = convertKeysToSnakeCase(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = convertKeysToSnakeCase(item)
		}
		return result
	default:
		return v
	}
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// unmarshalRequest unmarshals JSON request into proto message
func (h *HTTPJSONHandler) unmarshalRequest(r *http.Request, msg proto.Message) error {
	unmarshaler := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	decoder := json.NewDecoder(r.Body)
	var rawJSON map[string]interface{}
	if err := decoder.Decode(&rawJSON); err != nil {
		return err
	}

	// Convert camelCase keys to snake_case for proto compatibility
	// The TypeScript client sends camelCase (deckList) but proto expects snake_case (deck_list)
	convertedJSON := convertKeysToSnakeCase(rawJSON).(map[string]interface{})

	jsonBytes, err := json.Marshal(convertedJSON)
	if err != nil {
		return err
	}

	return unmarshaler.Unmarshal(jsonBytes, msg)
}

// unmarshalRequestWithDebug unmarshals JSON request into proto message and returns raw JSON for debugging
func (h *HTTPJSONHandler) unmarshalRequestWithDebug(r *http.Request, msg proto.Message) (map[string]interface{}, error) {
	unmarshaler := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	decoder := json.NewDecoder(r.Body)
	var rawJSON map[string]interface{}
	if err := decoder.Decode(&rawJSON); err != nil {
		return nil, err
	}

	// Convert camelCase keys to snake_case for proto compatibility
	convertedJSON := convertKeysToSnakeCase(rawJSON).(map[string]interface{})

	jsonBytes, err := json.Marshal(convertedJSON)
	if err != nil {
		return rawJSON, err
	}

	return rawJSON, unmarshaler.Unmarshal(jsonBytes, msg)
}

// writeSuccessResponse writes a successful JSON response
func (h *HTTPJSONHandler) writeSuccessResponse(w http.ResponseWriter, msg proto.Message) {
	marshaler := protojson.MarshalOptions{
		EmitUnpopulated: true,
		UseProtoNames:   true,
	}

	jsonBytes, err := marshaler.Marshal(msg)
	if err != nil {
		h.logger.Error("failed to marshal response", zap.Error(err))
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

// writeErrorResponse writes an error JSON response
func (h *HTTPJSONHandler) writeErrorResponse(w http.ResponseWriter, err error) {
	h.logger.Error("RPC error", zap.Error(err))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	errorResp := map[string]string{
		"error": err.Error(),
	}

	json.NewEncoder(w).Encode(errorResp)
}

// ==================== Additional Handler Implementations ====================

// handleConnectAdmin handles the ConnectAdmin method
func (h *HTTPJSONHandler) handleConnectAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ConnectAdminRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ConnectAdmin(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleConnectSetUserData handles the ConnectSetUserData method
func (h *HTTPJSONHandler) handleConnectSetUserData(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ConnectSetUserDataRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ConnectSetUserData(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAuthSendTokenToEmail handles the AuthSendTokenToEmail method
func (h *HTTPJSONHandler) handleAuthSendTokenToEmail(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AuthSendTokenToEmailRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AuthSendTokenToEmail(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAuthResetPassword handles the AuthResetPassword method
func (h *HTTPJSONHandler) handleAuthResetPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AuthResetPasswordRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AuthResetPassword(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleServerGetPromotionMessages handles the ServerGetPromotionMessages method
func (h *HTTPJSONHandler) handleServerGetPromotionMessages(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ServerGetPromotionMessagesRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ServerGetPromotionMessages(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleServerAddFeedbackMessage handles the ServerAddFeedbackMessage method
func (h *HTTPJSONHandler) handleServerAddFeedbackMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ServerAddFeedbackMessageRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ServerAddFeedbackMessage(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomGetUsers handles the RoomGetUsers method
func (h *HTTPJSONHandler) handleRoomGetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomGetUsersRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomGetUsers(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomGetFinishedMatches handles the RoomGetFinishedMatches method
func (h *HTTPJSONHandler) handleRoomGetFinishedMatches(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomGetFinishedMatchesRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomGetFinishedMatches(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomGetTableById handles the RoomGetTableById method
func (h *HTTPJSONHandler) handleRoomGetTableById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomGetTableByIdRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomGetTableById(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomCreateTable handles the RoomCreateTable method
func (h *HTTPJSONHandler) handleRoomCreateTable(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomCreateTableRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomCreateTable(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomCreateTournament handles the RoomCreateTournament method
func (h *HTTPJSONHandler) handleRoomCreateTournament(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomCreateTournamentRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomCreateTournament(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomJoinTable handles the RoomJoinTable method
func (h *HTTPJSONHandler) handleRoomJoinTable(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomJoinTableRequest
	rawJSON, err := h.unmarshalRequestWithDebug(r, &req)
	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Debug logging to trace deck submission - show raw JSON keys
	rawKeys := make([]string, 0)
	for k := range rawJSON {
		rawKeys = append(rawKeys, k)
	}
	deckListRaw := ""
	if v, ok := rawJSON["deckList"]; ok {
		if s, ok := v.(string); ok {
			deckListRaw = s
		}
	}
	if v, ok := rawJSON["deck_list"]; ok {
		if s, ok := v.(string); ok {
			deckListRaw = s
		}
	}

	h.logger.Info("[HTTP DEBUG] RoomJoinTable raw request",
		zap.Strings("raw_json_keys", rawKeys),
		zap.Int("raw_deck_list_length", len(deckListRaw)),
		zap.String("raw_deck_preview", truncateString(deckListRaw, 200)),
	)

	h.logger.Info("[HTTP DEBUG] RoomJoinTable request parsed",
		zap.String("table_id", req.GetTableId()),
		zap.Int("deck_list_length", len(req.GetDeckList())),
		zap.Bool("has_deck", req.GetDeckList() != ""),
	)

	resp, err := h.mageServer.RoomJoinTable(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomJoinTournament handles the RoomJoinTournament method
func (h *HTTPJSONHandler) handleRoomJoinTournament(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomJoinTournamentRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomJoinTournament(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomLeaveTableOrTournament handles the RoomLeaveTableOrTournament method
func (h *HTTPJSONHandler) handleRoomLeaveTableOrTournament(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomLeaveTableOrTournamentRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomLeaveTableOrTournament(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomWatchTable handles the RoomWatchTable method
func (h *HTTPJSONHandler) handleRoomWatchTable(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomWatchTableRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomWatchTable(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleRoomWatchTournament handles the RoomWatchTournament method
func (h *HTTPJSONHandler) handleRoomWatchTournament(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.RoomWatchTournamentRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.RoomWatchTournament(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleTableSwapSeats handles the TableSwapSeats method
func (h *HTTPJSONHandler) handleTableSwapSeats(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.TableSwapSeatsRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.TableSwapSeats(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleTableRemove handles the TableRemove method
func (h *HTTPJSONHandler) handleTableRemove(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.TableRemoveRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.TableRemove(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleTableIsOwner handles the TableIsOwner method
func (h *HTTPJSONHandler) handleTableIsOwner(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.TableIsOwnerRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.TableIsOwner(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDeckSubmit handles the DeckSubmit method
func (h *HTTPJSONHandler) handleDeckSubmit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DeckSubmitRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DeckSubmit(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleGameJoin handles the GameJoin method
func (h *HTTPJSONHandler) handleGameJoin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handleGameJoin started")

	var req pb.GameJoinRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		h.logger.Error("handleGameJoin unmarshal error", zap.Error(err))
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("handleGameJoin calling mageServer.GameJoin",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
	)

	resp, err := h.mageServer.GameJoin(ctx, &req)
	if err != nil {
		h.logger.Error("handleGameJoin error", zap.Error(err))
		h.writeErrorResponse(w, err)
		return
	}

	h.logger.Info("handleGameJoin success",
		zap.String("game_id", req.GetGameId()),
		zap.Bool("success", resp.GetSuccess()),
	)
	h.writeSuccessResponse(w, resp)
}

// handleGameWatchStart handles the GameWatchStart method
func (h *HTTPJSONHandler) handleGameWatchStart(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.GameWatchStartRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.GameWatchStart(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleGameWatchStop handles the GameWatchStop method
func (h *HTTPJSONHandler) handleGameWatchStop(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.GameWatchStopRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.GameWatchStop(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleGameGetView handles the GameGetView method
func (h *HTTPJSONHandler) handleGameGetView(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.logger.Info("handleGameGetView started")

	var req pb.GameGetViewRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		h.logger.Error("handleGameGetView unmarshal error", zap.Error(err))
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Info("handleGameGetView calling mageServer.GameGetView",
		zap.String("session_id", req.GetSessionId()),
		zap.String("game_id", req.GetGameId()),
		zap.String("player_id", req.GetPlayerId()),
	)

	resp, err := h.mageServer.GameGetView(ctx, &req)
	if err != nil {
		h.logger.Error("handleGameGetView error", zap.Error(err))
		h.writeErrorResponse(w, err)
		return
	}

	h.logger.Info("handleGameGetView success",
		zap.String("game_id", req.GetGameId()),
		zap.Bool("has_game", resp.GetGame() != nil),
	)
	h.writeSuccessResponse(w, resp)
}

// handleSendPlayerUUID handles the SendPlayerUUID method
func (h *HTTPJSONHandler) handleSendPlayerUUID(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendPlayerUUIDRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendPlayerUUID(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendPlayerString handles the SendPlayerString method
func (h *HTTPJSONHandler) handleSendPlayerString(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendPlayerStringRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendPlayerString(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendPlayerBoolean handles the SendPlayerBoolean method
func (h *HTTPJSONHandler) handleSendPlayerBoolean(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendPlayerBooleanRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendPlayerBoolean(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendPlayerInteger handles the SendPlayerInteger method
func (h *HTTPJSONHandler) handleSendPlayerInteger(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendPlayerIntegerRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendPlayerInteger(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendPlayerManaType handles the SendPlayerManaType method
func (h *HTTPJSONHandler) handleSendPlayerManaType(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendPlayerManaTypeRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendPlayerManaType(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendPlayerAction handles the SendPlayerAction method
func (h *HTTPJSONHandler) handleSendPlayerAction(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendPlayerActionRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendPlayerAction(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleMatchStart handles the MatchStart method
func (h *HTTPJSONHandler) handleMatchStart(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.MatchStartRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.MatchStart(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleMatchQuit handles the MatchQuit method
func (h *HTTPJSONHandler) handleMatchQuit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.MatchQuitRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.MatchQuit(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendSpecialAction handles the SendSpecialAction method (play land, foretell, etc.)
func (h *HTTPJSONHandler) handleSendSpecialAction(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendSpecialActionRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendSpecialAction(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDraftJoin handles the DraftJoin method
func (h *HTTPJSONHandler) handleDraftJoin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DraftJoinRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DraftJoin(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendDraftCardPick handles the SendDraftCardPick method
func (h *HTTPJSONHandler) handleSendDraftCardPick(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendDraftCardPickRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendDraftCardPick(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleSendDraftCardMark handles the SendDraftCardMark method
func (h *HTTPJSONHandler) handleSendDraftCardMark(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.SendDraftCardMarkRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.SendDraftCardMark(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDraftSetBoosterLoaded handles the DraftSetBoosterLoaded method
func (h *HTTPJSONHandler) handleDraftSetBoosterLoaded(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DraftSetBoosterLoadedRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DraftSetBoosterLoaded(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleDraftQuit handles the DraftQuit method
func (h *HTTPJSONHandler) handleDraftQuit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.DraftQuitRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.DraftQuit(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleTournamentJoin handles the TournamentJoin method
func (h *HTTPJSONHandler) handleTournamentJoin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.TournamentJoinRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.TournamentJoin(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleTournamentStart handles the TournamentStart method
func (h *HTTPJSONHandler) handleTournamentStart(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.TournamentStartRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.TournamentStart(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleTournamentQuit handles the TournamentQuit method
func (h *HTTPJSONHandler) handleTournamentQuit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.TournamentQuitRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.TournamentQuit(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleTournamentFindById handles the TournamentFindById method
func (h *HTTPJSONHandler) handleTournamentFindById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.TournamentFindByIdRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.TournamentFindById(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleChatJoin handles the ChatJoin method
func (h *HTTPJSONHandler) handleChatJoin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ChatJoinRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ChatJoin(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleChatLeave handles the ChatLeave method
func (h *HTTPJSONHandler) handleChatLeave(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ChatLeaveRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ChatLeave(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleChatSendMessage handles the ChatSendMessage method
func (h *HTTPJSONHandler) handleChatSendMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ChatSendMessageRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ChatSendMessage(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleChatFindByTable handles the ChatFindByTable method
func (h *HTTPJSONHandler) handleChatFindByTable(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ChatFindByTableRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ChatFindByTable(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleChatFindByGame handles the ChatFindByGame method
func (h *HTTPJSONHandler) handleChatFindByGame(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ChatFindByGameRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ChatFindByGame(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleChatFindByTournament handles the ChatFindByTournament method
func (h *HTTPJSONHandler) handleChatFindByTournament(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ChatFindByTournamentRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ChatFindByTournament(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleChatFindByRoom handles the ChatFindByRoom method
func (h *HTTPJSONHandler) handleChatFindByRoom(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ChatFindByRoomRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ChatFindByRoom(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleGetMatchHistory handles the GetMatchHistory method
func (h *HTTPJSONHandler) handleGetMatchHistory(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.GetMatchHistoryRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.GetMatchHistory(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleGetMatchById handles the GetMatchById method
func (h *HTTPJSONHandler) handleGetMatchById(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.GetMatchByIdRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.GetMatchById(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleReplayInit handles the ReplayInit method
func (h *HTTPJSONHandler) handleReplayInit(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ReplayInitRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ReplayInit(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleReplayStart handles the ReplayStart method
func (h *HTTPJSONHandler) handleReplayStart(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ReplayStartRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ReplayStart(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleReplayStop handles the ReplayStop method
func (h *HTTPJSONHandler) handleReplayStop(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ReplayStopRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ReplayStop(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleReplayNext handles the ReplayNext method
func (h *HTTPJSONHandler) handleReplayNext(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ReplayNextRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ReplayNext(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleReplayPrevious handles the ReplayPrevious method
func (h *HTTPJSONHandler) handleReplayPrevious(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ReplayPreviousRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ReplayPrevious(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleReplaySkipForward handles the ReplaySkipForward method
func (h *HTTPJSONHandler) handleReplaySkipForward(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.ReplaySkipForwardRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.ReplaySkipForward(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminGetUsers handles the AdminGetUsers method
func (h *HTTPJSONHandler) handleAdminGetUsers(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminGetUsersRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminGetUsers(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminDisconnectUser handles the AdminDisconnectUser method
func (h *HTTPJSONHandler) handleAdminDisconnectUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminDisconnectUserRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminDisconnectUser(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminMuteUser handles the AdminMuteUser method
func (h *HTTPJSONHandler) handleAdminMuteUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminMuteUserRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminMuteUser(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminLockUser handles the AdminLockUser method
func (h *HTTPJSONHandler) handleAdminLockUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminLockUserRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminLockUser(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminActivateUser handles the AdminActivateUser method
func (h *HTTPJSONHandler) handleAdminActivateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminActivateUserRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminActivateUser(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminToggleActivateUser handles the AdminToggleActivateUser method
func (h *HTTPJSONHandler) handleAdminToggleActivateUser(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminToggleActivateUserRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminToggleActivateUser(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminEndUserSession handles the AdminEndUserSession method
func (h *HTTPJSONHandler) handleAdminEndUserSession(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminEndUserSessionRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminEndUserSession(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminTableRemove handles the AdminTableRemove method
func (h *HTTPJSONHandler) handleAdminTableRemove(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminTableRemoveRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminTableRemove(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}

// handleAdminSendBroadcastMessage handles the AdminSendBroadcastMessage method
func (h *HTTPJSONHandler) handleAdminSendBroadcastMessage(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var req pb.AdminSendBroadcastMessageRequest
	if err := h.unmarshalRequest(r, &req); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.mageServer.AdminSendBroadcastMessage(ctx, &req)
	if err != nil {
		h.writeErrorResponse(w, err)
		return
	}

	h.writeSuccessResponse(w, resp)
}
