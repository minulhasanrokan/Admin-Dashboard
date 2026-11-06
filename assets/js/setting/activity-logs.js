'use strict';

(function ($) {
    let activityGrid = null;

    $(document).ready(function () {
        initActivityLogsGrid();
    });

    function initActivityLogsGrid() {
        activityGrid = initAgDataGrid({
            gridContainer: '#activity_ag_grid',
            url: '/settings/activity-logs/list',
            exportFileName: 'system_activity_logs',
            defaultOrderCol: 'created_at',
            defaultOrderDir: 'desc',
            totalBadgeSelector: '#activity_total_badge',
            columnDefs: [
                {
                    headerName: '#',
                    isIndexCol: true,
                    width: 70,
                    minWidth: 65
                },
                {
                    field: 'activity',
                    colId: 'activity',
                    headerName: 'Activity Description',
                    width: 260,
                    minWidth: 200,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="fw-medium text-dark"><i class="ti ti-activity me-1 text-primary"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'activity_status',
                    colId: 'activity_status',
                    headerName: 'Action Type',
                    width: 140,
                    minWidth: 120,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        let status = params.value != null ? String(params.value).toLowerCase(): '';
                        let badgeClass = 'bg-secondary-transparent text-secondary';
                        if (status.includes('create') || status.includes('add')) {
                            badgeClass = 'bg-success-transparent text-success';
                        } else if (status.includes('edit') || status.includes('update')) {
                            badgeClass = 'bg-info-transparent text-info';
                        } else if (status.includes('delete') || status.includes('remove')) {
                            badgeClass = 'bg-danger-transparent text-danger';
                        }
                        return `<span class="badge ${badgeClass} fw-semibold">${params.value || 'Action'}</span>`;
                    }
                },
                {
                    field: 'table_name',
                    colId: 'table_name',
                    headerName: 'Target Entity',
                    width: 160,
                    minWidth: 140,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="badge bg-light text-dark border">${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'ip_address',
                    colId: 'ip_address',
                    headerName: 'IP Address',
                    width: 150,
                    minWidth: 130,
                    sortable: true,
                    filter: true,
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-network me-1"></i>${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'device',
                    colId: 'device',
                    headerName: 'Device / User Agent',
                    width: 220,
                    minWidth: 180,
                    sortable: false,
                    cellRenderer: function (params) {
                        return `<span class="text-truncate d-inline-block mw-100 fs-12 text-muted" title="${params.value || ''}">${params.value || '-'}</span>`;
                    }
                },
                {
                    field: 'created_at',
                    colId: 'created_at',
                    headerName: 'Timestamp',
                    width: 200,
                    minWidth: 180,
                    sortable: true,
                    filter: true,
                    filterType: 'date',
                    cellRenderer: function (params) {
                        return `<span class="text-muted fs-12"><i class="ti ti-clock me-1"></i>${params.value || '-'}</span>`;
                    }
                }
            ]
        });
    }
})(jQuery);
