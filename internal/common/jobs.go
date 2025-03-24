package common

// ====================================================
// DATASTRUCTURES FOR JOB REQUESTS AND RESPONSE
// ====================================================

type Job interface {
	Execute() JobS
}
