package admin

// boundedListCap is the safety ceiling for list queries that are NOT exposed
// to clients with offset/limit pagination. These lists describe intrinsically
// small sets (areas per org, courses per org, subjects per area, students per
// course, time slots per course). The cap is far above any realistic tenant
// usage — it exists so a misconfigured seed or runaway insert cannot make the
// API return thousands of rows per request.
//
// For lists that can grow unboundedly and must be paginated, see
// providers.Pagination (topic/activity endpoints).
const boundedListCap = 500
