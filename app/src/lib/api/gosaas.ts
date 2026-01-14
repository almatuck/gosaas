import webapi from "./gocliRequest"
import * as components from "./gosaasComponents"
export * from "./gosaasComponents"

/**
 * @description "Health check endpoint"
 */
export function healthCheck() {
	return webapi.get<components.HealthResponse>(`/health`)
}

/**
 * @description "Request password reset"
 * @param req
 */
export function forgotPassword(req: components.ForgotPasswordRequest) {
	return webapi.post<components.MessageResponse>(`/api/v1/auth/forgot-password`, req)
}

/**
 * @description "User login"
 * @param req
 */
export function login(req: components.LoginRequest) {
	return webapi.post<components.LoginResponse>(`/api/v1/auth/login`, req)
}

/**
 * @description "Refresh authentication token"
 * @param req
 */
export function refreshToken(req: components.RefreshTokenRequest) {
	return webapi.post<components.RefreshTokenResponse>(`/api/v1/auth/refresh`, req)
}

/**
 * @description "Register new user"
 * @param req
 */
export function register(req: components.RegisterRequest) {
	return webapi.post<components.LoginResponse>(`/api/v1/auth/register`, req)
}

/**
 * @description "Resend email verification"
 * @param req
 */
export function resendVerification(req: components.ResendVerificationRequest) {
	return webapi.post<components.MessageResponse>(`/api/v1/auth/resend-verification`, req)
}

/**
 * @description "Reset password with token"
 * @param req
 */
export function resetPassword(req: components.ResetPasswordRequest) {
	return webapi.post<components.MessageResponse>(`/api/v1/auth/reset-password`, req)
}

/**
 * @description "Verify email address with token"
 * @param req
 */
export function verifyEmail(req: components.EmailVerificationRequest) {
	return webapi.post<components.MessageResponse>(`/api/v1/auth/verify-email`, req)
}

/**
 * @description "List all available subscription plans"
 */
export function listPlans() {
	return webapi.get<components.ListPlansResponse>(`/api/v1/subscription/plans`)
}

/**
 * @description "Get current user's subscription"
 */
export function getSubscription() {
	return webapi.get<components.GetSubscriptionResponse>(`/api/v1/subscription`)
}

/**
 * @description "Get billing history"
 * @param params
 */
export function listBillingHistory(params: components.ListBillingHistoryRequestParams) {
	return webapi.get<components.ListBillingHistoryResponse>(`/api/v1/subscription/billing-history`, params)
}

/**
 * @description "Create billing portal session for subscription management"
 */
export function createBillingPortal() {
	return webapi.post<components.CreateBillingPortalResponse>(`/api/v1/subscription/billing-portal`)
}

/**
 * @description "Cancel current subscription"
 * @param req
 */
export function cancelSubscription(req: components.CancelSubscriptionRequest) {
	return webapi.post<components.CancelSubscriptionResponse>(`/api/v1/subscription/cancel`, req)
}

/**
 * @description "Check if user has access to a feature"
 * @param req
 */
export function checkFeature(req: components.CheckFeatureRequest) {
	return webapi.post<components.CheckFeatureResponse>(`/api/v1/subscription/check-feature`, req)
}

/**
 * @description "Create checkout session to subscribe to a plan"
 * @param req
 */
export function createCheckout(req: components.CreateCheckoutRequest) {
	return webapi.post<components.CreateCheckoutResponse>(`/api/v1/subscription/checkout`, req)
}

/**
 * @description "Get current usage statistics"
 */
export function getUsageStats() {
	return webapi.get<components.GetUsageStatsResponse>(`/api/v1/subscription/usage`)
}

/**
 * @description "Get current user profile"
 */
export function getCurrentUser() {
	return webapi.get<components.GetUserResponse>(`/api/v1/user/me`)
}

/**
 * @description "Update current user profile"
 * @param req
 */
export function updateCurrentUser(req: components.UpdateUserRequest) {
	return webapi.put<components.GetUserResponse>(`/api/v1/user/me`, req)
}

/**
 * @description "Delete current user account"
 * @param req
 */
export function deleteAccount(req: components.DeleteAccountRequest) {
	return webapi.delete<components.MessageResponse>(`/api/v1/user/me`, req)
}

/**
 * @description "Change password for authenticated user"
 * @param req
 */
export function changePassword(req: components.ChangePasswordRequest) {
	return webapi.post<components.MessageResponse>(`/api/v1/user/me/change-password`, req)
}

/**
 * @description "Get user preferences"
 */
export function getPreferences() {
	return webapi.get<components.GetPreferencesResponse>(`/api/v1/user/me/preferences`)
}

/**
 * @description "Update user preferences"
 * @param req
 */
export function updatePreferences(req: components.UpdatePreferencesRequest) {
	return webapi.put<components.GetPreferencesResponse>(`/api/v1/user/me/preferences`, req)
}
