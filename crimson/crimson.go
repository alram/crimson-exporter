package crimson

import (
	"bytes"
	"encoding/json"
	"log"
	"maps"
	"path/filepath"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	adminSocketPath = "/var/run/ceph"
)

type DumpMetrics struct {
	Metrics any `json:"metrics,omitempty"`
}

type CrimsonCollector struct {
	LBAAllocExtents                              *prometheus.Desc
	LBAAllocExtentsIterNexts                     *prometheus.Desc
	AlienReceiveBatchQueueLength                 *prometheus.Desc
	AlienTotalReceivedMessages                   *prometheus.Desc
	AlienTotalSentMessages                       *prometheus.Desc
	BackgroundProcessIoBlockedCount              *prometheus.Desc
	BackgroundProcessIoBlockedCountClean         *prometheus.Desc
	BackgroundProcessIoBlockedCountTrim          *prometheus.Desc
	BackgroundProcessIoBlockedSum                *prometheus.Desc
	BackgroundProcessIoCount                     *prometheus.Desc
	IoQueueActivations                           *prometheus.Desc
	IoQueueAdjustedConsumption                   *prometheus.Desc
	IoQueueConsumption                           *prometheus.Desc
	IoQueueDelay                                 *prometheus.Desc
	IoQueueDiskQueueLength                       *prometheus.Desc
	IoQueueFlowRatio                             *prometheus.Desc
	IoQueueQueueLength                           *prometheus.Desc
	IoQueueShares                                *prometheus.Desc
	IoQueueStarvationTimeSec                     *prometheus.Desc
	IoQueueTotalBytes                            *prometheus.Desc
	IoQueueTotalDelaySec                         *prometheus.Desc
	IoQueueTotalExecSec                          *prometheus.Desc
	IoQueueTotalOperations                       *prometheus.Desc
	IoQueueTotalReadBytes                        *prometheus.Desc
	IoQueueTotalReadOps                          *prometheus.Desc
	IoQueueTotalSplitBytes                       *prometheus.Desc
	IoQueueTotalSplitOps                         *prometheus.Desc
	IoQueueTotalWriteBytes                       *prometheus.Desc
	IoQueueTotalWriteOps                         *prometheus.Desc
	JournalIoDepthNum                            *prometheus.Desc
	JournalIoNum                                 *prometheus.Desc
	JournalRecordGroupDataBytes                  *prometheus.Desc
	JournalRecordGroupMetadataBytes              *prometheus.Desc
	JournalRecordGroupPaddingBytes               *prometheus.Desc
	JournalRecordNum                             *prometheus.Desc
	JournalTrimmerAllocJournalBytes              *prometheus.Desc
	JournalTrimmerDirtyJournalBytes              *prometheus.Desc
	MemoryAllocatedMemory                        *prometheus.Desc
	MemoryCrossCPUFreeOperations                 *prometheus.Desc
	MemoryFreeMemory                             *prometheus.Desc
	MemoryFreeOperations                         *prometheus.Desc
	MemoryMallocFailed                           *prometheus.Desc
	MemoryMallocLiveObjects                      *prometheus.Desc
	MemoryMallocOperations                       *prometheus.Desc
	MemoryReclaimsOperations                     *prometheus.Desc
	MemoryTotalMemory                            *prometheus.Desc
	NetworkBytesReceived                         *prometheus.Desc
	NetworkBytesSent                             *prometheus.Desc
	ReactorAbandonedFailedFutures                *prometheus.Desc
	ReactorAioBytesRead                          *prometheus.Desc
	ReactorAioBytesWrite                         *prometheus.Desc
	ReactorAioErrors                             *prometheus.Desc
	ReactorAioOutsizes                           *prometheus.Desc
	ReactorAioReads                              *prometheus.Desc
	ReactorAioWrites                             *prometheus.Desc
	ReactorAwakeTimeMsTotal                      *prometheus.Desc
	ReactorCppExceptions                         *prometheus.Desc
	ReactorCPUBusyMs                             *prometheus.Desc
	ReactorCPUStealTimeMs                        *prometheus.Desc
	ReactorCPUUsedTimeMs                         *prometheus.Desc
	ReactorFstreamReadBytes                      *prometheus.Desc
	ReactorFstreamReadBytesBlocked               *prometheus.Desc
	ReactorFstreamReads                          *prometheus.Desc
	ReactorFstreamReadsAheadBytesDiscarded       *prometheus.Desc
	ReactorFstreamReadsAheadsDiscarded           *prometheus.Desc
	ReactorFstreamReadsBlocked                   *prometheus.Desc
	ReactorFsyncs                                *prometheus.Desc
	ReactorIoThreadedFallbacks                   *prometheus.Desc
	ReactorLoggingFailures                       *prometheus.Desc
	ReactorPolls                                 *prometheus.Desc
	ReactorSleepTimeMsTotal                      *prometheus.Desc
	ReactorStalls                                *prometheus.Desc
	ReactorTasksPending                          *prometheus.Desc
	ReactorTasksProcessed                        *prometheus.Desc
	ReactorTimersPending                         *prometheus.Desc
	ReactorUtilization                           *prometheus.Desc
	SchedulerQueueLength                         *prometheus.Desc
	SchedulerRuntimeMs                           *prometheus.Desc
	SchedulerShares                              *prometheus.Desc
	SchedulerStarvetimeMs                        *prometheus.Desc
	SchedulerTasksProcessed                      *prometheus.Desc
	SchedulerTimeSpentOnTaskQuotaViolationsMs    *prometheus.Desc
	SchedulerWaittimeMs                          *prometheus.Desc
	SeastoreConcurrentTransactions               *prometheus.Desc
	SeastoreOpLat                                *prometheus.Desc
	SeastorePendingTransactions                  *prometheus.Desc
	SegmentCleanerAvailableBytes                 *prometheus.Desc
	SegmentCleanerAvailableRatio                 *prometheus.Desc
	SegmentCleanerClosedJournalTotalBytes        *prometheus.Desc
	SegmentCleanerClosedJournalUsedBytes         *prometheus.Desc
	SegmentCleanerClosedOolTotalBytes            *prometheus.Desc
	SegmentCleanerClosedOolUsedBytes             *prometheus.Desc
	SegmentCleanerProjectedCount                 *prometheus.Desc
	SegmentCleanerProjectedUsedBytesSum          *prometheus.Desc
	SegmentCleanerReclaimRatio                   *prometheus.Desc
	SegmentCleanerReclaimedBytes                 *prometheus.Desc
	SegmentCleanerReclaimedSegmentBytes          *prometheus.Desc
	SegmentCleanerSegmentSize                    *prometheus.Desc
	SegmentCleanerSegmentUtilizationDistribution *prometheus.Desc
	SegmentCleanerSegmentsClosed                 *prometheus.Desc
	SegmentCleanerSegmentsCountCloseJournal      *prometheus.Desc
	SegmentCleanerSegmentsCountCloseOol          *prometheus.Desc
	SegmentCleanerSegmentsCountOpenJournal       *prometheus.Desc
	SegmentCleanerSegmentsCountOpenOol           *prometheus.Desc
	SegmentCleanerSegmentsCountReleaseJournal    *prometheus.Desc
	SegmentCleanerSegmentsCountReleaseOol        *prometheus.Desc
	SegmentCleanerSegmentsEmpty                  *prometheus.Desc
	SegmentCleanerSegmentsInJournal              *prometheus.Desc
	SegmentCleanerSegmentsNumber                 *prometheus.Desc
	SegmentCleanerSegmentsOpen                   *prometheus.Desc
	SegmentCleanerSegmentsTypeJournal            *prometheus.Desc
	SegmentCleanerSegmentsTypeOol                *prometheus.Desc
	SegmentCleanerTotalBytes                     *prometheus.Desc
	SegmentCleanerUnavailableReclaimableBytes    *prometheus.Desc
	SegmentCleanerUnavailableUnreclaimableBytes  *prometheus.Desc
	SegmentCleanerUnavailableUnusedBytes         *prometheus.Desc
	SegmentCleanerUsedBytes                      *prometheus.Desc
	SegmentManagerClosedSegments                 *prometheus.Desc
	SegmentManagerClosedSegmentsUnusedBytes      *prometheus.Desc
	SegmentManagerDataReadBytes                  *prometheus.Desc
	SegmentManagerDataReadNum                    *prometheus.Desc
	SegmentManagerDataWriteBytes                 *prometheus.Desc
	SegmentManagerDataWriteNum                   *prometheus.Desc
	SegmentManagerMetadataWriteBytes             *prometheus.Desc
	SegmentManagerMetadataWriteNum               *prometheus.Desc
	SegmentManagerOpenedSegments                 *prometheus.Desc
	SegmentManagerReleasedSegments               *prometheus.Desc
	StallDetectorReported                        *prometheus.Desc

	//in case i wanna add some tests, stub-able funcs
	getOSDSocketFiles  func(adminSocketPath string) ([]string, error)
	runAsokDumpMetrics func(socketFilename string) ([]byte, error)
	getOSDID           func(socketFilename string) (string, error)
}

type CrimsonData struct {
	OSDID     string
	NumShards int

	Name map[string][]struct {
		Shard       string
		Value       any //can be int or struct{}
		ExtraLabels map[string]string
	}
}

func NewCrimsonExporter() *CrimsonCollector {

	// Bunch of these labels are similar but I still want to separate them
	// in different labels vars since crimson is still under dev. if there are
	// any changes that'll just make life less painful to adapt.
	shardLabel := []string{"osd", "shard"}
	ioLabel := []string{"osd", "shard"}
	journalLabel := []string{"osd", "shard", "submitter"}
	networkLabel := []string{"osd", "group", "shard"}
	memoryLabel := []string{"osd", "shard"}
	reactorLabel := []string{"osd", "shard"}
	reactorStallLabel := []string{"osd", "shard", "latency", "le"}
	schedulerLabel := []string{"osd", "group", "shard"}
	seatoreLatencyLabel := []string{"osd", "shard", "le", "latency"}
	seastoreLabel := []string{"osd", "shard"}
	segmentLabel := []string{"osd", "shard"}
	segmentCleanerUtilizationLabel := []string{"osd", "shard", "le"}
	segmentManagerLabel := []string{"osd", "device_id", "shard"}

	return &CrimsonCollector{

		getOSDSocketFiles:  GetOSDSocketFiles,
		runAsokDumpMetrics: RunAsokDumpMetrics,
		getOSDID:           GetOSDID,

		LBAAllocExtents:                              prometheus.NewDesc("crimson_osd_lba_alloc_extents", "", shardLabel, nil),
		LBAAllocExtentsIterNexts:                     prometheus.NewDesc("crimson_osd_lba_alloc_extents_iter_nexts", "", shardLabel, nil),
		AlienReceiveBatchQueueLength:                 prometheus.NewDesc("crimson_osd_alien_receive_batch_queue_length", "", shardLabel, nil),
		AlienTotalReceivedMessages:                   prometheus.NewDesc("crimson_osd_alien_total_received_messages", "", shardLabel, nil),
		AlienTotalSentMessages:                       prometheus.NewDesc("crimson_osd_alien_total_sent_messages", "", shardLabel, nil),
		BackgroundProcessIoBlockedCount:              prometheus.NewDesc("crimson_osd_background_process_io_blocked_count", "", shardLabel, nil),
		BackgroundProcessIoBlockedCountClean:         prometheus.NewDesc("crimson_osd_background_process_io_blocked_count_clean", "", shardLabel, nil),
		BackgroundProcessIoBlockedCountTrim:          prometheus.NewDesc("crimson_osd_background_process_io_blocked_count_trim", "", shardLabel, nil),
		BackgroundProcessIoBlockedSum:                prometheus.NewDesc("crimson_osd_background_process_io_blocked_sum", "", shardLabel, nil),
		BackgroundProcessIoCount:                     prometheus.NewDesc("crimson_osd_background_process_io_count", "", shardLabel, nil),
		IoQueueActivations:                           prometheus.NewDesc("crimson_osd_io_queue_activations", "", ioLabel, nil),
		IoQueueAdjustedConsumption:                   prometheus.NewDesc("crimson_osd_io_queue_adjusted_consumption", "", ioLabel, nil),
		IoQueueConsumption:                           prometheus.NewDesc("crimson_osd_io_queue_consumption", "", ioLabel, nil),
		IoQueueDelay:                                 prometheus.NewDesc("crimson_osd_io_queue_delay", "", ioLabel, nil),
		IoQueueDiskQueueLength:                       prometheus.NewDesc("crimson_osd_io_queue_disk_queue_length", "", ioLabel, nil),
		IoQueueFlowRatio:                             prometheus.NewDesc("crimson_osd_io_queue_flow_ratio", "", ioLabel, nil),
		IoQueueQueueLength:                           prometheus.NewDesc("crimson_osd_io_queue_queue_length", "", ioLabel, nil),
		IoQueueShares:                                prometheus.NewDesc("crimson_osd_io_queue_shares", "", ioLabel, nil),
		IoQueueStarvationTimeSec:                     prometheus.NewDesc("crimson_osd_io_queue_starvation_time_sec", "", ioLabel, nil),
		IoQueueTotalBytes:                            prometheus.NewDesc("crimson_osd_io_queue_total_bytes", "", ioLabel, nil),
		IoQueueTotalDelaySec:                         prometheus.NewDesc("crimson_osd_io_queue_total_delay_sec", "", ioLabel, nil),
		IoQueueTotalExecSec:                          prometheus.NewDesc("crimson_osd_io_queue_total_exec_sec", "", ioLabel, nil),
		IoQueueTotalOperations:                       prometheus.NewDesc("crimson_osd_io_queue_total_operations", "", ioLabel, nil),
		IoQueueTotalReadBytes:                        prometheus.NewDesc("crimson_osd_io_queue_total_read_bytes", "", ioLabel, nil),
		IoQueueTotalReadOps:                          prometheus.NewDesc("crimson_osd_io_queue_total_read_ops", "", ioLabel, nil),
		IoQueueTotalSplitBytes:                       prometheus.NewDesc("crimson_osd_io_queue_total_split_bytes", "", ioLabel, nil),
		IoQueueTotalSplitOps:                         prometheus.NewDesc("crimson_osd_io_queue_total_split_ops", "", ioLabel, nil),
		IoQueueTotalWriteBytes:                       prometheus.NewDesc("crimson_osd_io_queue_total_write_bytes", "", ioLabel, nil),
		IoQueueTotalWriteOps:                         prometheus.NewDesc("crimson_osd_io_queue_total_write_ops", "", ioLabel, nil),
		JournalIoDepthNum:                            prometheus.NewDesc("crimson_osd_journal_io_depth_num", "", journalLabel, nil),
		JournalIoNum:                                 prometheus.NewDesc("crimson_osd_journal_io_num", "", journalLabel, nil),
		JournalRecordGroupDataBytes:                  prometheus.NewDesc("crimson_osd_journal_record_group_data_bytes", "", journalLabel, nil),
		JournalRecordGroupMetadataBytes:              prometheus.NewDesc("crimson_osd_journal_record_group_metadata_bytes", "", journalLabel, nil),
		JournalRecordGroupPaddingBytes:               prometheus.NewDesc("crimson_osd_journal_record_group_padding_bytes", "", journalLabel, nil),
		JournalRecordNum:                             prometheus.NewDesc("crimson_osd_journal_record_num", "", journalLabel, nil),
		JournalTrimmerAllocJournalBytes:              prometheus.NewDesc("crimson_osd_journal_trimmer_alloc_journal_bytes", "", journalLabel, nil),
		JournalTrimmerDirtyJournalBytes:              prometheus.NewDesc("crimson_osd_journal_trimmer_dirty_journal_bytes", "", journalLabel, nil),
		MemoryAllocatedMemory:                        prometheus.NewDesc("crimson_osd_memory_allocated_memory", "", memoryLabel, nil),
		MemoryCrossCPUFreeOperations:                 prometheus.NewDesc("crimson_osd_memory_cross_cpu_free_operations", "", memoryLabel, nil),
		MemoryFreeMemory:                             prometheus.NewDesc("crimson_osd_memory_free_memory", "", memoryLabel, nil),
		MemoryFreeOperations:                         prometheus.NewDesc("crimson_osd_memory_free_operations", "", memoryLabel, nil),
		MemoryMallocFailed:                           prometheus.NewDesc("crimson_osd_memory_malloc_failed", "", memoryLabel, nil),
		MemoryMallocLiveObjects:                      prometheus.NewDesc("crimson_osd_memory_malloc_live_objects", "", memoryLabel, nil),
		MemoryMallocOperations:                       prometheus.NewDesc("crimson_osd_memory_malloc_operations", "", memoryLabel, nil),
		MemoryReclaimsOperations:                     prometheus.NewDesc("crimson_osd_memory_reclaims_operations", "", memoryLabel, nil),
		MemoryTotalMemory:                            prometheus.NewDesc("crimson_osd_memory_total_memory", "", memoryLabel, nil),
		NetworkBytesReceived:                         prometheus.NewDesc("crimson_osd_network_bytes_received", "", networkLabel, nil),
		NetworkBytesSent:                             prometheus.NewDesc("crimson_osd_network_bytes_sent", "", networkLabel, nil),
		ReactorAbandonedFailedFutures:                prometheus.NewDesc("crimson_osd_reactor_abandoned_failed_futures", "", reactorLabel, nil),
		ReactorAioBytesRead:                          prometheus.NewDesc("crimson_osd_reactor_aio_bytes_read", "", reactorLabel, nil),
		ReactorAioBytesWrite:                         prometheus.NewDesc("crimson_osd_reactor_aio_bytes_write", "", reactorLabel, nil),
		ReactorAioErrors:                             prometheus.NewDesc("crimson_osd_reactor_aio_errors", "", reactorLabel, nil),
		ReactorAioOutsizes:                           prometheus.NewDesc("crimson_osd_reactor_aio_outsizes", "", reactorLabel, nil),
		ReactorAioReads:                              prometheus.NewDesc("crimson_osd_reactor_aio_reads", "", reactorLabel, nil),
		ReactorAioWrites:                             prometheus.NewDesc("crimson_osd_reactor_aio_writes", "", reactorLabel, nil),
		ReactorAwakeTimeMsTotal:                      prometheus.NewDesc("crimson_osd_reactor_awake_time_ms_total", "", reactorLabel, nil),
		ReactorCppExceptions:                         prometheus.NewDesc("crimson_osd_reactor_cpp_exceptions", "", reactorLabel, nil),
		ReactorCPUBusyMs:                             prometheus.NewDesc("crimson_osd_reactor_cpu_busy_ms", "", reactorLabel, nil),
		ReactorCPUStealTimeMs:                        prometheus.NewDesc("crimson_osd_reactor_cpu_steal_time_ms", "", reactorLabel, nil),
		ReactorCPUUsedTimeMs:                         prometheus.NewDesc("crimson_osd_reactor_cpu_used_time_ms", "", reactorLabel, nil),
		ReactorFstreamReadBytes:                      prometheus.NewDesc("crimson_osd_reactor_fstream_read_bytes", "", reactorLabel, nil),
		ReactorFstreamReadBytesBlocked:               prometheus.NewDesc("crimson_osd_reactor_fstream_read_bytes_blocked", "", reactorLabel, nil),
		ReactorFstreamReads:                          prometheus.NewDesc("crimson_osd_reactor_fstream_reads", "", reactorLabel, nil),
		ReactorFstreamReadsAheadBytesDiscarded:       prometheus.NewDesc("crimson_osd_reactor_fstream_reads_ahead_bytes_discarded", "", reactorLabel, nil),
		ReactorFstreamReadsAheadsDiscarded:           prometheus.NewDesc("crimson_osd_reactor_fstream_reads_aheads_discarded", "", reactorLabel, nil),
		ReactorFstreamReadsBlocked:                   prometheus.NewDesc("crimson_osd_reactor_fstream_reads_blocked", "", reactorLabel, nil),
		ReactorFsyncs:                                prometheus.NewDesc("crimson_osd_reactor_fsyncs", "", reactorLabel, nil),
		ReactorIoThreadedFallbacks:                   prometheus.NewDesc("crimson_osd_reactor_io_threaded_fallbacks", "", reactorLabel, nil),
		ReactorLoggingFailures:                       prometheus.NewDesc("crimson_osd_reactor_logging_failures", "", reactorLabel, nil),
		ReactorPolls:                                 prometheus.NewDesc("crimson_osd_reactor_polls", "", reactorLabel, nil),
		ReactorSleepTimeMsTotal:                      prometheus.NewDesc("crimson_osd_reactor_sleep_time_ms_total", "", reactorLabel, nil),
		ReactorStalls:                                prometheus.NewDesc("crimson_osd_reactor_stalls", "", reactorStallLabel, nil),
		ReactorTasksPending:                          prometheus.NewDesc("crimson_osd_reactor_tasks_pending", "", reactorLabel, nil),
		ReactorTasksProcessed:                        prometheus.NewDesc("crimson_osd_reactor_tasks_processed", "", reactorLabel, nil),
		ReactorTimersPending:                         prometheus.NewDesc("crimson_osd_reactor_timers_pending", "", reactorLabel, nil),
		ReactorUtilization:                           prometheus.NewDesc("crimson_osd_reactor_utilization", "", reactorLabel, nil),
		SchedulerQueueLength:                         prometheus.NewDesc("crimson_osd_scheduler_queue_length", "", schedulerLabel, nil),
		SchedulerRuntimeMs:                           prometheus.NewDesc("crimson_osd_scheduler_runtime_ms", "", schedulerLabel, nil),
		SchedulerShares:                              prometheus.NewDesc("crimson_osd_scheduler_shares", "", schedulerLabel, nil),
		SchedulerStarvetimeMs:                        prometheus.NewDesc("crimson_osd_scheduler_starvetime_ms", "", schedulerLabel, nil),
		SchedulerTasksProcessed:                      prometheus.NewDesc("crimson_osd_scheduler_tasks_processed", "", schedulerLabel, nil),
		SchedulerTimeSpentOnTaskQuotaViolationsMs:    prometheus.NewDesc("crimson_osd_scheduler_time_spent_on_task_quota_violations_ms", "", schedulerLabel, nil),
		SchedulerWaittimeMs:                          prometheus.NewDesc("crimson_osd_scheduler_waittime_ms", "", schedulerLabel, nil),
		SeastoreConcurrentTransactions:               prometheus.NewDesc("crimson_osd_seastore_concurrent_transactions", "", seastoreLabel, nil),
		SeastoreOpLat:                                prometheus.NewDesc("crimson_osd_seastore_op_lat", "", seatoreLatencyLabel, nil),
		SeastorePendingTransactions:                  prometheus.NewDesc("crimson_osd_seastore_pending_transactions", "", seastoreLabel, nil),
		SegmentCleanerAvailableBytes:                 prometheus.NewDesc("crimson_osd_segment_cleaner_available_bytes", "", segmentLabel, nil),
		SegmentCleanerAvailableRatio:                 prometheus.NewDesc("crimson_osd_segment_cleaner_available_ratio", "", segmentLabel, nil),
		SegmentCleanerClosedJournalTotalBytes:        prometheus.NewDesc("crimson_osd_segment_cleaner_closed_journal_total_bytes", "", segmentLabel, nil),
		SegmentCleanerClosedJournalUsedBytes:         prometheus.NewDesc("crimson_osd_segment_cleaner_closed_journal_used_bytes", "", segmentLabel, nil),
		SegmentCleanerClosedOolTotalBytes:            prometheus.NewDesc("crimson_osd_segment_cleaner_closed_ool_total_bytes", "", segmentLabel, nil),
		SegmentCleanerClosedOolUsedBytes:             prometheus.NewDesc("crimson_osd_segment_cleaner_closed_ool_used_bytes", "", segmentLabel, nil),
		SegmentCleanerProjectedCount:                 prometheus.NewDesc("crimson_osd_segment_cleaner_projected_count", "", segmentLabel, nil),
		SegmentCleanerProjectedUsedBytesSum:          prometheus.NewDesc("crimson_osd_segment_cleaner_projected_used_bytes_sum", "", segmentLabel, nil),
		SegmentCleanerReclaimRatio:                   prometheus.NewDesc("crimson_osd_segment_cleaner_reclaim_ratio", "", segmentLabel, nil),
		SegmentCleanerReclaimedBytes:                 prometheus.NewDesc("crimson_osd_segment_cleaner_reclaimed_bytes", "", segmentLabel, nil),
		SegmentCleanerReclaimedSegmentBytes:          prometheus.NewDesc("crimson_osd_segment_cleaner_reclaimed_segment_bytes", "", segmentLabel, nil),
		SegmentCleanerSegmentSize:                    prometheus.NewDesc("crimson_osd_segment_cleaner_segment_size", "", segmentLabel, nil),
		SegmentCleanerSegmentUtilizationDistribution: prometheus.NewDesc("crimson_osd_segment_cleaner_segment_utilization_distribution", "", segmentCleanerUtilizationLabel, nil),
		SegmentCleanerSegmentsClosed:                 prometheus.NewDesc("crimson_osd_segment_cleaner_segments_closed", "", segmentLabel, nil),
		SegmentCleanerSegmentsCountCloseJournal:      prometheus.NewDesc("crimson_osd_segment_cleaner_segments_count_close_journal", "", segmentLabel, nil),
		SegmentCleanerSegmentsCountCloseOol:          prometheus.NewDesc("crimson_osd_segment_cleaner_segments_count_close_ool", "", segmentLabel, nil),
		SegmentCleanerSegmentsCountOpenJournal:       prometheus.NewDesc("crimson_osd_segment_cleaner_segments_count_open_journal", "", segmentLabel, nil),
		SegmentCleanerSegmentsCountOpenOol:           prometheus.NewDesc("crimson_osd_segment_cleaner_segments_count_open_ool", "", segmentLabel, nil),
		SegmentCleanerSegmentsCountReleaseJournal:    prometheus.NewDesc("crimson_osd_segment_cleaner_segments_count_release_journal", "", segmentLabel, nil),
		SegmentCleanerSegmentsCountReleaseOol:        prometheus.NewDesc("crimson_osd_segment_cleaner_segments_count_release_ool", "", segmentLabel, nil),
		SegmentCleanerSegmentsEmpty:                  prometheus.NewDesc("crimson_osd_segment_cleaner_segments_empty", "", segmentLabel, nil),
		SegmentCleanerSegmentsInJournal:              prometheus.NewDesc("crimson_osd_segment_cleaner_segments_in_journal", "", segmentLabel, nil),
		SegmentCleanerSegmentsNumber:                 prometheus.NewDesc("crimson_osd_segment_cleaner_segments_number", "", segmentLabel, nil),
		SegmentCleanerSegmentsOpen:                   prometheus.NewDesc("crimson_osd_segment_cleaner_segments_open", "", segmentLabel, nil),
		SegmentCleanerSegmentsTypeJournal:            prometheus.NewDesc("crimson_osd_segment_cleaner_segments_type_journal", "", segmentLabel, nil),
		SegmentCleanerSegmentsTypeOol:                prometheus.NewDesc("crimson_osd_segment_cleaner_segments_type_ool", "", segmentLabel, nil),
		SegmentCleanerTotalBytes:                     prometheus.NewDesc("crimson_osd_segment_cleaner_total_bytes", "", segmentLabel, nil),
		SegmentCleanerUnavailableReclaimableBytes:    prometheus.NewDesc("crimson_osd_segment_cleaner_unavailable_reclaimable_bytes", "", segmentLabel, nil),
		SegmentCleanerUnavailableUnreclaimableBytes:  prometheus.NewDesc("crimson_osd_segment_cleaner_unavailable_unreclaimable_bytes", "", segmentLabel, nil),
		SegmentCleanerUnavailableUnusedBytes:         prometheus.NewDesc("crimson_osd_segment_cleaner_unavailable_unused_bytes", "", segmentLabel, nil),
		SegmentCleanerUsedBytes:                      prometheus.NewDesc("crimson_osd_segment_cleaner_used_bytes", "", segmentLabel, nil),
		SegmentManagerClosedSegments:                 prometheus.NewDesc("crimson_osd_segment_manager_closed_segments", "", segmentManagerLabel, nil),
		SegmentManagerClosedSegmentsUnusedBytes:      prometheus.NewDesc("crimson_osd_segment_manager_closed_segments_unused_bytes", "", segmentManagerLabel, nil),
		SegmentManagerDataReadBytes:                  prometheus.NewDesc("crimson_osd_segment_manager_data_read_bytes", "", segmentManagerLabel, nil),
		SegmentManagerDataReadNum:                    prometheus.NewDesc("crimson_osd_segment_manager_data_read_num", "", segmentManagerLabel, nil),
		SegmentManagerDataWriteBytes:                 prometheus.NewDesc("crimson_osd_segment_manager_data_write_bytes", "", segmentManagerLabel, nil),
		SegmentManagerDataWriteNum:                   prometheus.NewDesc("crimson_osd_segment_manager_data_write_num", "", segmentManagerLabel, nil),
		SegmentManagerMetadataWriteBytes:             prometheus.NewDesc("crimson_osd_segment_manager_metadata_write_bytes", "", segmentManagerLabel, nil),
		SegmentManagerMetadataWriteNum:               prometheus.NewDesc("crimson_osd_segment_manager_metadata_write_num", "", segmentManagerLabel, nil),
		SegmentManagerOpenedSegments:                 prometheus.NewDesc("crimson_osd_segment_manager_opened_segments", "", segmentManagerLabel, nil),
		SegmentManagerReleasedSegments:               prometheus.NewDesc("crimson_osd_segment_manager_released_segments", "", segmentManagerLabel, nil),
		StallDetectorReported:                        prometheus.NewDesc("crimson_osd_stall_detector_reported", "", shardLabel, nil),
	}
}

func (c *CrimsonCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.LBAAllocExtents
	ch <- c.LBAAllocExtentsIterNexts
	ch <- c.AlienReceiveBatchQueueLength
	ch <- c.AlienTotalReceivedMessages
	ch <- c.AlienTotalSentMessages
	ch <- c.BackgroundProcessIoBlockedCount
	ch <- c.BackgroundProcessIoBlockedCountClean
	ch <- c.BackgroundProcessIoBlockedCountTrim
	ch <- c.BackgroundProcessIoBlockedSum
	ch <- c.BackgroundProcessIoCount
	ch <- c.IoQueueActivations
	ch <- c.IoQueueAdjustedConsumption
	ch <- c.IoQueueConsumption
	ch <- c.IoQueueDelay
	ch <- c.IoQueueDiskQueueLength
	ch <- c.IoQueueFlowRatio
	ch <- c.IoQueueQueueLength
	ch <- c.IoQueueShares
	ch <- c.IoQueueStarvationTimeSec
	ch <- c.IoQueueTotalBytes
	ch <- c.IoQueueTotalDelaySec
	ch <- c.IoQueueTotalExecSec
	ch <- c.IoQueueTotalOperations
	ch <- c.IoQueueTotalReadBytes
	ch <- c.IoQueueTotalReadOps
	ch <- c.IoQueueTotalSplitBytes
	ch <- c.IoQueueTotalSplitOps
	ch <- c.IoQueueTotalWriteBytes
	ch <- c.IoQueueTotalWriteOps
	ch <- c.JournalIoDepthNum
	ch <- c.JournalIoNum
	ch <- c.JournalRecordGroupDataBytes
	ch <- c.JournalRecordGroupMetadataBytes
	ch <- c.JournalRecordGroupPaddingBytes
	ch <- c.JournalRecordNum
	ch <- c.JournalTrimmerAllocJournalBytes
	ch <- c.JournalTrimmerDirtyJournalBytes
	ch <- c.MemoryAllocatedMemory
	ch <- c.MemoryCrossCPUFreeOperations
	ch <- c.MemoryFreeMemory
	ch <- c.MemoryFreeOperations
	ch <- c.MemoryMallocFailed
	ch <- c.MemoryMallocLiveObjects
	ch <- c.MemoryMallocOperations
	ch <- c.MemoryReclaimsOperations
	ch <- c.MemoryTotalMemory
	ch <- c.NetworkBytesReceived
	ch <- c.NetworkBytesSent
	ch <- c.ReactorAbandonedFailedFutures
	ch <- c.ReactorAioBytesRead
	ch <- c.ReactorAioBytesWrite
	ch <- c.ReactorAioErrors
	ch <- c.ReactorAioOutsizes
	ch <- c.ReactorAioReads
	ch <- c.ReactorAioWrites
	ch <- c.ReactorAwakeTimeMsTotal
	ch <- c.ReactorCppExceptions
	ch <- c.ReactorCPUBusyMs
	ch <- c.ReactorCPUStealTimeMs
	ch <- c.ReactorCPUUsedTimeMs
	ch <- c.ReactorFstreamReadBytes
	ch <- c.ReactorFstreamReadBytesBlocked
	ch <- c.ReactorFstreamReads
	ch <- c.ReactorFstreamReadsAheadBytesDiscarded
	ch <- c.ReactorFstreamReadsAheadsDiscarded
	ch <- c.ReactorFstreamReadsBlocked
	ch <- c.ReactorFsyncs
	ch <- c.ReactorIoThreadedFallbacks
	ch <- c.ReactorLoggingFailures
	ch <- c.ReactorPolls
	ch <- c.ReactorSleepTimeMsTotal
	ch <- c.ReactorStalls
	ch <- c.ReactorTasksPending
	ch <- c.ReactorTasksProcessed
	ch <- c.ReactorTimersPending
	ch <- c.ReactorUtilization
	ch <- c.SchedulerQueueLength
	ch <- c.SchedulerRuntimeMs
	ch <- c.SchedulerShares
	ch <- c.SchedulerStarvetimeMs
	ch <- c.SchedulerTasksProcessed
	ch <- c.SchedulerTimeSpentOnTaskQuotaViolationsMs
	ch <- c.SchedulerWaittimeMs
	ch <- c.SeastoreConcurrentTransactions
	ch <- c.SeastoreOpLat
	ch <- c.SeastorePendingTransactions
	ch <- c.SegmentCleanerAvailableBytes
	ch <- c.SegmentCleanerAvailableRatio
	ch <- c.SegmentCleanerClosedJournalTotalBytes
	ch <- c.SegmentCleanerClosedJournalUsedBytes
	ch <- c.SegmentCleanerClosedOolTotalBytes
	ch <- c.SegmentCleanerClosedOolUsedBytes
	ch <- c.SegmentCleanerProjectedCount
	ch <- c.SegmentCleanerProjectedUsedBytesSum
	ch <- c.SegmentCleanerReclaimRatio
	ch <- c.SegmentCleanerReclaimedBytes
	ch <- c.SegmentCleanerReclaimedSegmentBytes
	ch <- c.SegmentCleanerSegmentSize
	ch <- c.SegmentCleanerSegmentUtilizationDistribution
	ch <- c.SegmentCleanerSegmentsClosed
	ch <- c.SegmentCleanerSegmentsCountCloseJournal
	ch <- c.SegmentCleanerSegmentsCountCloseOol
	ch <- c.SegmentCleanerSegmentsCountOpenJournal
	ch <- c.SegmentCleanerSegmentsCountOpenOol
	ch <- c.SegmentCleanerSegmentsCountReleaseJournal
	ch <- c.SegmentCleanerSegmentsCountReleaseOol
	ch <- c.SegmentCleanerSegmentsEmpty
	ch <- c.SegmentCleanerSegmentsInJournal
	ch <- c.SegmentCleanerSegmentsNumber
	ch <- c.SegmentCleanerSegmentsOpen
	ch <- c.SegmentCleanerSegmentsTypeJournal
	ch <- c.SegmentCleanerSegmentsTypeOol
	ch <- c.SegmentCleanerTotalBytes
	ch <- c.SegmentCleanerUnavailableReclaimableBytes
	ch <- c.SegmentCleanerUnavailableUnreclaimableBytes
	ch <- c.SegmentCleanerUnavailableUnusedBytes
	ch <- c.SegmentCleanerUsedBytes
	ch <- c.SegmentManagerClosedSegments
	ch <- c.SegmentManagerClosedSegmentsUnusedBytes
	ch <- c.SegmentManagerDataReadBytes
	ch <- c.SegmentManagerDataReadNum
	ch <- c.SegmentManagerDataWriteBytes
	ch <- c.SegmentManagerDataWriteNum
	ch <- c.SegmentManagerMetadataWriteBytes
	ch <- c.SegmentManagerMetadataWriteNum
	ch <- c.SegmentManagerOpenedSegments
	ch <- c.SegmentManagerReleasedSegments
	ch <- c.StallDetectorReported
}

func (c *CrimsonCollector) Collect(ch chan<- prometheus.Metric) {
	metrics, err := c.getAdminStats()
	if err != nil {
		log.Fatalf("failed to get crimson metrics: %v", err)
	}

	var wg sync.WaitGroup
	for _, metric := range metrics {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < metric.NumShards; i++ {
				ch <- prometheus.MustNewConstMetric(c.LBAAllocExtents, prometheus.GaugeValue, metric.Name["LBA_alloc_extents"][i].Value.(float64), metric.OSDID, metric.Name["LBA_alloc_extents"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.LBAAllocExtentsIterNexts, prometheus.GaugeValue, metric.Name["LBA_alloc_extents_iter_nexts"][i].Value.(float64), metric.OSDID, metric.Name["LBA_alloc_extents_iter_nexts"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.AlienReceiveBatchQueueLength, prometheus.GaugeValue, metric.Name["alien_receive_batch_queue_length"][i].Value.(float64), metric.OSDID, metric.Name["alien_receive_batch_queue_length"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.AlienTotalReceivedMessages, prometheus.GaugeValue, metric.Name["alien_total_received_messages"][i].Value.(float64), metric.OSDID, metric.Name["alien_total_received_messages"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.AlienTotalSentMessages, prometheus.GaugeValue, metric.Name["alien_total_sent_messages"][i].Value.(float64), metric.OSDID, metric.Name["alien_total_sent_messages"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.BackgroundProcessIoBlockedCount, prometheus.GaugeValue, metric.Name["background_process_io_blocked_count"][i].Value.(float64), metric.OSDID, metric.Name["background_process_io_blocked_count"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.BackgroundProcessIoBlockedCountClean, prometheus.GaugeValue, metric.Name["background_process_io_blocked_count_clean"][i].Value.(float64), metric.OSDID, metric.Name["background_process_io_blocked_count_clean"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.BackgroundProcessIoBlockedCountTrim, prometheus.GaugeValue, metric.Name["background_process_io_blocked_count_trim"][i].Value.(float64), metric.OSDID, metric.Name["background_process_io_blocked_count_trim"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.BackgroundProcessIoBlockedSum, prometheus.GaugeValue, metric.Name["background_process_io_blocked_sum"][i].Value.(float64), metric.OSDID, metric.Name["background_process_io_blocked_sum"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.BackgroundProcessIoCount, prometheus.GaugeValue, metric.Name["background_process_io_count"][i].Value.(float64), metric.OSDID, metric.Name["background_process_io_count"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueActivations, prometheus.GaugeValue, metric.Name["io_queue_activations"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_activations"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueAdjustedConsumption, prometheus.GaugeValue, metric.Name["io_queue_adjusted_consumption"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_adjusted_consumption"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueConsumption, prometheus.GaugeValue, metric.Name["io_queue_consumption"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_consumption"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueDelay, prometheus.GaugeValue, metric.Name["io_queue_delay"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_delay"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueDiskQueueLength, prometheus.GaugeValue, metric.Name["io_queue_disk_queue_length"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_disk_queue_length"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueFlowRatio, prometheus.GaugeValue, metric.Name["io_queue_flow_ratio"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_flow_ratio"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueQueueLength, prometheus.GaugeValue, metric.Name["io_queue_queue_length"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_queue_length"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueShares, prometheus.GaugeValue, metric.Name["io_queue_shares"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_shares"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueStarvationTimeSec, prometheus.GaugeValue, metric.Name["io_queue_starvation_time_sec"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_starvation_time_sec"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalBytes, prometheus.GaugeValue, metric.Name["io_queue_total_bytes"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalDelaySec, prometheus.GaugeValue, metric.Name["io_queue_total_delay_sec"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_delay_sec"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalExecSec, prometheus.GaugeValue, metric.Name["io_queue_total_exec_sec"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_exec_sec"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalOperations, prometheus.GaugeValue, metric.Name["io_queue_total_operations"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_operations"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalReadBytes, prometheus.GaugeValue, metric.Name["io_queue_total_read_bytes"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_read_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalReadOps, prometheus.GaugeValue, metric.Name["io_queue_total_read_ops"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_read_ops"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalSplitBytes, prometheus.GaugeValue, metric.Name["io_queue_total_split_bytes"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_split_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalSplitOps, prometheus.GaugeValue, metric.Name["io_queue_total_split_ops"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_split_ops"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalWriteBytes, prometheus.GaugeValue, metric.Name["io_queue_total_write_bytes"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_write_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.IoQueueTotalWriteOps, prometheus.GaugeValue, metric.Name["io_queue_total_write_ops"][i].Value.(float64), metric.OSDID, metric.Name["io_queue_total_write_ops"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.JournalIoDepthNum, prometheus.GaugeValue, metric.Name["journal_io_depth_num"][i].Value.(float64), metric.OSDID, metric.Name["journal_io_depth_num"][i].Shard, metric.Name["journal_io_depth_num"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.JournalIoNum, prometheus.GaugeValue, metric.Name["journal_io_num"][i].Value.(float64), metric.OSDID, metric.Name["journal_io_num"][i].Shard, metric.Name["journal_io_num"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.JournalRecordGroupDataBytes, prometheus.GaugeValue, metric.Name["journal_record_group_data_bytes"][i].Value.(float64), metric.OSDID, metric.Name["journal_record_group_data_bytes"][i].Shard, metric.Name["journal_record_group_data_bytes"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.JournalRecordGroupMetadataBytes, prometheus.GaugeValue, metric.Name["journal_record_group_metadata_bytes"][i].Value.(float64), metric.OSDID, metric.Name["journal_record_group_metadata_bytes"][i].Shard, metric.Name["journal_record_group_metadata_bytes"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.JournalRecordGroupPaddingBytes, prometheus.GaugeValue, metric.Name["journal_record_group_padding_bytes"][i].Value.(float64), metric.OSDID, metric.Name["journal_record_group_padding_bytes"][i].Shard, metric.Name["journal_record_group_padding_bytes"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.JournalRecordNum, prometheus.GaugeValue, metric.Name["journal_record_num"][i].Value.(float64), metric.OSDID, metric.Name["journal_record_num"][i].Shard, metric.Name["journal_record_num"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.JournalTrimmerAllocJournalBytes, prometheus.GaugeValue, metric.Name["journal_trimmer_alloc_journal_bytes"][i].Value.(float64), metric.OSDID, metric.Name["journal_trimmer_alloc_journal_bytes"][i].Shard, metric.Name["journal_trimmer_alloc_journal_bytes"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.JournalTrimmerDirtyJournalBytes, prometheus.GaugeValue, metric.Name["journal_trimmer_dirty_journal_bytes"][i].Value.(float64), metric.OSDID, metric.Name["journal_trimmer_dirty_journal_bytes"][i].Shard, metric.Name["journal_trimmer_dirty_journal_bytes"][i].ExtraLabels["submitter"])
				ch <- prometheus.MustNewConstMetric(c.MemoryAllocatedMemory, prometheus.GaugeValue, metric.Name["memory_allocated_memory"][i].Value.(float64), metric.OSDID, metric.Name["memory_allocated_memory"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryCrossCPUFreeOperations, prometheus.GaugeValue, metric.Name["memory_cross_cpu_free_operations"][i].Value.(float64), metric.OSDID, metric.Name["memory_cross_cpu_free_operations"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryFreeMemory, prometheus.GaugeValue, metric.Name["memory_free_memory"][i].Value.(float64), metric.OSDID, metric.Name["memory_free_memory"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryFreeOperations, prometheus.GaugeValue, metric.Name["memory_free_operations"][i].Value.(float64), metric.OSDID, metric.Name["memory_free_operations"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryMallocFailed, prometheus.GaugeValue, metric.Name["memory_malloc_failed"][i].Value.(float64), metric.OSDID, metric.Name["memory_malloc_failed"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryMallocLiveObjects, prometheus.GaugeValue, metric.Name["memory_malloc_live_objects"][i].Value.(float64), metric.OSDID, metric.Name["memory_malloc_live_objects"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryMallocOperations, prometheus.GaugeValue, metric.Name["memory_malloc_operations"][i].Value.(float64), metric.OSDID, metric.Name["memory_malloc_operations"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryReclaimsOperations, prometheus.GaugeValue, metric.Name["memory_reclaims_operations"][i].Value.(float64), metric.OSDID, metric.Name["memory_reclaims_operations"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.MemoryTotalMemory, prometheus.GaugeValue, metric.Name["memory_total_memory"][i].Value.(float64), metric.OSDID, metric.Name["memory_total_memory"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.NetworkBytesReceived, prometheus.GaugeValue, metric.Name["network_bytes_received"][i].Value.(float64), metric.OSDID, metric.Name["network_bytes_received"][i].Shard, metric.Name["network_bytes_received"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.NetworkBytesSent, prometheus.GaugeValue, metric.Name["network_bytes_sent"][i].Value.(float64), metric.OSDID, metric.Name["network_bytes_sent"][i].Shard, metric.Name["network_bytes_sent"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.ReactorAbandonedFailedFutures, prometheus.GaugeValue, metric.Name["reactor_abandoned_failed_futures"][i].Value.(float64), metric.OSDID, metric.Name["reactor_abandoned_failed_futures"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorAioBytesRead, prometheus.GaugeValue, metric.Name["reactor_aio_bytes_read"][i].Value.(float64), metric.OSDID, metric.Name["reactor_aio_bytes_read"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorAioBytesWrite, prometheus.GaugeValue, metric.Name["reactor_aio_bytes_write"][i].Value.(float64), metric.OSDID, metric.Name["reactor_aio_bytes_write"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorAioErrors, prometheus.GaugeValue, metric.Name["reactor_aio_errors"][i].Value.(float64), metric.OSDID, metric.Name["reactor_aio_errors"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorAioOutsizes, prometheus.GaugeValue, metric.Name["reactor_aio_outsizes"][i].Value.(float64), metric.OSDID, metric.Name["reactor_aio_outsizes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorAioReads, prometheus.GaugeValue, metric.Name["reactor_aio_reads"][i].Value.(float64), metric.OSDID, metric.Name["reactor_aio_reads"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorAioWrites, prometheus.GaugeValue, metric.Name["reactor_aio_writes"][i].Value.(float64), metric.OSDID, metric.Name["reactor_aio_writes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorAwakeTimeMsTotal, prometheus.GaugeValue, metric.Name["reactor_awake_time_ms_total"][i].Value.(float64), metric.OSDID, metric.Name["reactor_awake_time_ms_total"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorCppExceptions, prometheus.GaugeValue, metric.Name["reactor_cpp_exceptions"][i].Value.(float64), metric.OSDID, metric.Name["reactor_cpp_exceptions"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorCPUBusyMs, prometheus.GaugeValue, metric.Name["reactor_cpu_busy_ms"][i].Value.(float64), metric.OSDID, metric.Name["reactor_cpu_busy_ms"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorCPUStealTimeMs, prometheus.GaugeValue, metric.Name["reactor_cpu_steal_time_ms"][i].Value.(float64), metric.OSDID, metric.Name["reactor_cpu_steal_time_ms"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorCPUUsedTimeMs, prometheus.GaugeValue, metric.Name["reactor_cpu_used_time_ms"][i].Value.(float64), metric.OSDID, metric.Name["reactor_cpu_used_time_ms"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorFstreamReadBytes, prometheus.GaugeValue, metric.Name["reactor_fstream_read_bytes"][i].Value.(float64), metric.OSDID, metric.Name["reactor_fstream_read_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorFstreamReadBytesBlocked, prometheus.GaugeValue, metric.Name["reactor_fstream_read_bytes_blocked"][i].Value.(float64), metric.OSDID, metric.Name["reactor_fstream_read_bytes_blocked"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorFstreamReads, prometheus.GaugeValue, metric.Name["reactor_fstream_reads"][i].Value.(float64), metric.OSDID, metric.Name["reactor_fstream_reads"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorFstreamReadsAheadBytesDiscarded, prometheus.GaugeValue, metric.Name["reactor_fstream_reads_ahead_bytes_discarded"][i].Value.(float64), metric.OSDID, metric.Name["reactor_fstream_reads_ahead_bytes_discarded"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorFstreamReadsAheadsDiscarded, prometheus.GaugeValue, metric.Name["reactor_fstream_reads_aheads_discarded"][i].Value.(float64), metric.OSDID, metric.Name["reactor_fstream_reads_aheads_discarded"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorFstreamReadsBlocked, prometheus.GaugeValue, metric.Name["reactor_fstream_reads_blocked"][i].Value.(float64), metric.OSDID, metric.Name["reactor_fstream_reads_blocked"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorFsyncs, prometheus.GaugeValue, metric.Name["reactor_fsyncs"][i].Value.(float64), metric.OSDID, metric.Name["reactor_fsyncs"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorIoThreadedFallbacks, prometheus.GaugeValue, metric.Name["reactor_io_threaded_fallbacks"][i].Value.(float64), metric.OSDID, metric.Name["reactor_io_threaded_fallbacks"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorLoggingFailures, prometheus.GaugeValue, metric.Name["reactor_logging_failures"][i].Value.(float64), metric.OSDID, metric.Name["reactor_logging_failures"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorPolls, prometheus.GaugeValue, metric.Name["reactor_polls"][i].Value.(float64), metric.OSDID, metric.Name["reactor_polls"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorSleepTimeMsTotal, prometheus.GaugeValue, metric.Name["reactor_sleep_time_ms_total"][i].Value.(float64), metric.OSDID, metric.Name["reactor_sleep_time_ms_total"][i].Shard)
				// This is a latency bucket metric, needs some care
				//ch <- prometheus.MustNewConstMetric(c.ReactorStalls, prometheus.GaugeValue, metric.Name["reactor_stalls"][i].Value.(float64), metric.OSDID, metric.Name["reactor_stalls"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorTasksPending, prometheus.GaugeValue, metric.Name["reactor_tasks_pending"][i].Value.(float64), metric.OSDID, metric.Name["reactor_tasks_pending"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorTasksProcessed, prometheus.GaugeValue, metric.Name["reactor_tasks_processed"][i].Value.(float64), metric.OSDID, metric.Name["reactor_tasks_processed"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorTimersPending, prometheus.GaugeValue, metric.Name["reactor_timers_pending"][i].Value.(float64), metric.OSDID, metric.Name["reactor_timers_pending"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.ReactorUtilization, prometheus.GaugeValue, metric.Name["reactor_utilization"][i].Value.(float64), metric.OSDID, metric.Name["reactor_utilization"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SchedulerQueueLength, prometheus.GaugeValue, metric.Name["scheduler_queue_length"][i].Value.(float64), metric.OSDID, metric.Name["scheduler_queue_length"][i].Shard, metric.Name["scheduler_queue_length"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.SchedulerRuntimeMs, prometheus.GaugeValue, metric.Name["scheduler_runtime_ms"][i].Value.(float64), metric.OSDID, metric.Name["scheduler_runtime_ms"][i].Shard, metric.Name["scheduler_runtime_ms"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.SchedulerShares, prometheus.GaugeValue, metric.Name["scheduler_shares"][i].Value.(float64), metric.OSDID, metric.Name["scheduler_shares"][i].Shard, metric.Name["scheduler_shares"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.SchedulerStarvetimeMs, prometheus.GaugeValue, metric.Name["scheduler_starvetime_ms"][i].Value.(float64), metric.OSDID, metric.Name["scheduler_starvetime_ms"][i].Shard, metric.Name["scheduler_starvetime_ms"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.SchedulerTasksProcessed, prometheus.GaugeValue, metric.Name["scheduler_tasks_processed"][i].Value.(float64), metric.OSDID, metric.Name["scheduler_tasks_processed"][i].Shard, metric.Name["scheduler_tasks_processed"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.SchedulerTimeSpentOnTaskQuotaViolationsMs, prometheus.GaugeValue, metric.Name["scheduler_time_spent_on_task_quota_violations_ms"][i].Value.(float64), metric.OSDID, metric.Name["scheduler_time_spent_on_task_quota_violations_ms"][i].Shard, metric.Name["scheduler_time_spent_on_task_quota_violations_ms"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.SchedulerWaittimeMs, prometheus.GaugeValue, metric.Name["scheduler_waittime_ms"][i].Value.(float64), metric.OSDID, metric.Name["scheduler_waittime_ms"][i].Shard, metric.Name["scheduler_waittime_ms"][i].ExtraLabels["group"])
				ch <- prometheus.MustNewConstMetric(c.SeastoreConcurrentTransactions, prometheus.GaugeValue, metric.Name["seastore_concurrent_transactions"][i].Value.(float64), metric.OSDID, metric.Name["seastore_concurrent_transactions"][i].Shard)
				// This is a latency bucket metric, needs some care
				//ch <- prometheus.MustNewConstMetric(c.SeastoreOpLat, prometheus.GaugeValue, metric.Name["seastore_op_lat"][i].Value.(float64), metric.OSDID, metric.Name["seastore_op_lat"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SeastorePendingTransactions, prometheus.GaugeValue, metric.Name["seastore_pending_transactions"][i].Value.(float64), metric.OSDID, metric.Name["seastore_pending_transactions"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerAvailableBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_available_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_available_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerAvailableRatio, prometheus.GaugeValue, metric.Name["segment_cleaner_available_ratio"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_available_ratio"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerClosedJournalTotalBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_closed_journal_total_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_closed_journal_total_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerClosedJournalUsedBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_closed_journal_used_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_closed_journal_used_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerClosedOolTotalBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_closed_ool_total_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_closed_ool_total_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerClosedOolUsedBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_closed_ool_used_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_closed_ool_used_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerProjectedCount, prometheus.GaugeValue, metric.Name["segment_cleaner_projected_count"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_projected_count"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerProjectedUsedBytesSum, prometheus.GaugeValue, metric.Name["segment_cleaner_projected_used_bytes_sum"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_projected_used_bytes_sum"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerReclaimRatio, prometheus.GaugeValue, metric.Name["segment_cleaner_reclaim_ratio"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_reclaim_ratio"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerReclaimedBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_reclaimed_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_reclaimed_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerReclaimedSegmentBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_reclaimed_segment_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_reclaimed_segment_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentSize, prometheus.GaugeValue, metric.Name["segment_cleaner_segment_size"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segment_size"][i].Shard)
				// This is a latency bucket metric, needs some care
				//ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentUtilizationDistribution, prometheus.GaugeValue, metric.Name["segment_cleaner_segment_utilization_distribution"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segment_utilization_distribution"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsClosed, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_closed"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_closed"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsCountCloseJournal, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_count_close_journal"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_count_close_journal"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsCountCloseOol, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_count_close_ool"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_count_close_ool"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsCountOpenJournal, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_count_open_journal"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_count_open_journal"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsCountOpenOol, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_count_open_ool"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_count_open_ool"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsCountReleaseJournal, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_count_release_journal"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_count_release_journal"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsCountReleaseOol, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_count_release_ool"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_count_release_ool"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsEmpty, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_empty"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_empty"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsInJournal, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_in_journal"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_in_journal"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsNumber, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_number"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_number"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsOpen, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_open"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_open"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsTypeJournal, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_type_journal"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_type_journal"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerSegmentsTypeOol, prometheus.GaugeValue, metric.Name["segment_cleaner_segments_type_ool"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_segments_type_ool"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerTotalBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_total_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_total_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerUnavailableReclaimableBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_unavailable_reclaimable_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_unavailable_reclaimable_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerUnavailableUnreclaimableBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_unavailable_unreclaimable_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_unavailable_unreclaimable_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerUnavailableUnusedBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_unavailable_unused_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_unavailable_unused_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentCleanerUsedBytes, prometheus.GaugeValue, metric.Name["segment_cleaner_used_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_cleaner_used_bytes"][i].Shard)
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerClosedSegments, prometheus.GaugeValue, metric.Name["segment_manager_closed_segments"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_closed_segments"][i].Shard, metric.Name["segment_manager_closed_segments"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerClosedSegmentsUnusedBytes, prometheus.GaugeValue, metric.Name["segment_manager_closed_segments_unused_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_closed_segments_unused_bytes"][i].Shard, metric.Name["segment_manager_closed_segments_unused_bytes"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerDataReadBytes, prometheus.GaugeValue, metric.Name["segment_manager_data_read_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_data_read_bytes"][i].Shard, metric.Name["segment_manager_data_read_bytes"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerDataReadNum, prometheus.GaugeValue, metric.Name["segment_manager_data_read_num"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_data_read_num"][i].Shard, metric.Name["segment_manager_data_read_num"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerDataWriteBytes, prometheus.GaugeValue, metric.Name["segment_manager_data_write_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_data_write_bytes"][i].Shard, metric.Name["segment_manager_data_write_bytes"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerDataWriteNum, prometheus.GaugeValue, metric.Name["segment_manager_data_write_num"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_data_write_num"][i].Shard, metric.Name["segment_manager_data_write_num"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerMetadataWriteBytes, prometheus.GaugeValue, metric.Name["segment_manager_metadata_write_bytes"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_metadata_write_bytes"][i].Shard, metric.Name["segment_manager_metadata_write_bytes"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerMetadataWriteNum, prometheus.GaugeValue, metric.Name["segment_manager_metadata_write_num"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_metadata_write_num"][i].Shard, metric.Name["segment_manager_metadata_write_num"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerOpenedSegments, prometheus.GaugeValue, metric.Name["segment_manager_opened_segments"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_opened_segments"][i].Shard, metric.Name["segment_manager_opened_segments"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.SegmentManagerReleasedSegments, prometheus.GaugeValue, metric.Name["segment_manager_released_segments"][i].Value.(float64), metric.OSDID, metric.Name["segment_manager_released_segments"][i].Shard, metric.Name["segment_manager_released_segments"][i].ExtraLabels["device_id"])
				ch <- prometheus.MustNewConstMetric(c.StallDetectorReported, prometheus.GaugeValue, metric.Name["stall_detector_reported"][i].Value.(float64), metric.OSDID, metric.Name["stall_detector_reported"][i].Shard)
			}

		}()
	}
	wg.Wait()

}

func (c *CrimsonCollector) getAdminStats() ([]CrimsonData, error) {
	files, err := c.getOSDSocketFiles(adminSocketPath)
	if err != nil {
		return nil, err
	}
	var metrics []CrimsonData
	var lock sync.Mutex
	var wg sync.WaitGroup

	for _, f := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()

			osdID, err := c.getOSDID(f)
			if err != nil {
				log.Fatalf("failed to get osd id from socket %s: %v", f, err)
			}

			blob, err := c.runAsokDumpMetrics(filepath.Join(adminSocketPath, f))
			if err != nil {
				log.Fatal("failed extracting admin stats")
			}
			if blob == nil {
				return
			}

			// "+Inf" is not a valid integer, who would have thunk?
			blob = bytes.ReplaceAll(blob, []byte("\": \"+Inf\""), []byte("\": -1"))
			var metric DumpMetrics
			if err := json.Unmarshal(blob, &metric); err != nil {
				log.Fatalf("failed unmarshalling json: %v", err)
			}

			// We wanna transform the dumpmetrics map[string]any into something more consumable
			var crimsonMetric CrimsonData
			crimsonMetric.OSDID = osdID
			crimsonMetric.NumShards = 4
			crimsonMetric.Name = make(map[string][]struct {
				Shard       string
				Value       any
				ExtraLabels map[string]string
			})

			for _, m := range metric.Metrics.([]any) {
				for k := range maps.Keys(m.(map[string]any)) {
					shard := m.(map[string]any)[k].(map[string]any)["shard"].(string)
					val := m.(map[string]any)[k].(map[string]any)["value"]
					extraLabelValues := make(map[string]string)
					extraLabels := []string{
						"submitter",
						"group",
						"latency",
						"le",
						"device_id",
					}
					for _, label := range extraLabels {
						if _, ok := m.(map[string]any)[k].(map[string]any)[label]; ok {
							extraLabelValues[label] = m.(map[string]any)[k].(map[string]any)[label].(string)
						}
					}

					crimsonMetric.Name[k] = append(crimsonMetric.Name[k], struct {
						Shard       string
						Value       any
						ExtraLabels map[string]string
					}{
						Shard:       shard,
						Value:       val,
						ExtraLabels: extraLabelValues,
					})
				}
			}

			lock.Lock()
			metrics = append(metrics, crimsonMetric)
			lock.Unlock()
		}(f)
		wg.Wait()
	}
	return metrics, nil

}
