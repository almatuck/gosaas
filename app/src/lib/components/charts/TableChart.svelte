<script lang="ts">
	import { onMount } from 'svelte';
	import {
		ChevronUp,
		ChevronDown,
		ChevronsUpDown,
		Search,
		Download,
		Settings2,
		ChevronLeft,
		ChevronRight,
		ChevronsLeft,
		ChevronsRight,
		X,
		Check,
		Filter,
		ArrowUpDown
	} from 'lucide-svelte';
	import type { TableColumn, TableRow } from './types';

	interface Props {
		columns: TableColumn[];
		rows: TableRow[];
		height?: number;
		striped?: boolean;
		hoverable?: boolean;
		paginate?: boolean;
		pageSize?: number;
		searchable?: boolean;
		exportable?: boolean;
		selectable?: boolean;
		resizable?: boolean;
		showColumnToggle?: boolean;
		stickyHeader?: boolean;
		compactMode?: boolean;
		onRowClick?: (row: TableRow, index: number) => void;
		onSelectionChange?: (selectedRows: TableRow[]) => void;
	}

	let {
		columns: initialColumns = [],
		rows = [],
		height,
		striped = true,
		hoverable = true,
		paginate = true,
		pageSize: initialPageSize = 10,
		searchable = true,
		exportable = true,
		selectable = false,
		resizable = true,
		showColumnToggle = true,
		stickyHeader = true,
		compactMode = false,
		onRowClick,
		onSelectionChange
	}: Props = $props();

	// State
	let sortKey = $state<string | null>(null);
	let sortDirection = $state<'asc' | 'desc'>('asc');
	let currentPage = $state(0);
	let pageSize = $state(10);
	let searchQuery = $state('');
	let columnFilters = $state<Record<string, string>>({});
	let visibleColumns = $state<Set<string>>(new Set());
	let selectedRows = $state<Set<number>>(new Set());

	// Sync state with props when they change
	$effect(() => {
		pageSize = initialPageSize;
	});

	$effect(() => {
		visibleColumns = new Set(initialColumns.map((c) => c.key));
	});
	let columnWidths = $state<Record<string, number>>({});
	let showColumnMenu = $state(false);
	let showFilterMenu = $state<string | null>(null);
	let resizingColumn = $state<string | null>(null);
	let resizeStartX = $state(0);
	let resizeStartWidth = $state(0);

	// Refs
	let containerRef: HTMLDivElement;
	let tableRef: HTMLTableElement;

	// Computed columns with visibility
	const columns = $derived(initialColumns.filter((c) => visibleColumns.has(c.key)));

	// Filter and search rows
	const filteredRows = $derived.by(() => {
		let result = [...rows];

		// Apply global search
		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			result = result.filter((row) =>
				columns.some((col) => {
					const val = row[col.key];
					return val !== null && val !== undefined && String(val).toLowerCase().includes(query);
				})
			);
		}

		// Apply column filters
		for (const [key, filterValue] of Object.entries(columnFilters)) {
			if (filterValue.trim()) {
				const query = filterValue.toLowerCase();
				result = result.filter((row) => {
					const val = row[key];
					return val !== null && val !== undefined && String(val).toLowerCase().includes(query);
				});
			}
		}

		return result;
	});

	// Sort rows
	const sortedRows = $derived.by(() => {
		let sorted = [...filteredRows];

		if (sortKey) {
			const key = sortKey; // Capture for type narrowing
			sorted.sort((a, b) => {
				const aVal = a[key];
				const bVal = b[key];

				if (aVal === bVal) return 0;
				if (aVal === null || aVal === undefined) return 1;
				if (bVal === null || bVal === undefined) return -1;

				let comparison = 0;
				if (typeof aVal === 'number' && typeof bVal === 'number') {
					comparison = aVal - bVal;
				} else if (aVal instanceof Date && bVal instanceof Date) {
					comparison = aVal.getTime() - bVal.getTime();
				} else {
					comparison = String(aVal).localeCompare(String(bVal));
				}

				return sortDirection === 'asc' ? comparison : -comparison;
			});
		}

		return sorted;
	});

	// Paginate rows
	const paginatedRows = $derived.by(() => {
		if (!paginate) return sortedRows;
		const start = currentPage * pageSize;
		const end = start + pageSize;
		return sortedRows.slice(start, end);
	});

	const totalPages = $derived(Math.ceil(filteredRows.length / pageSize));
	const totalRows = $derived(filteredRows.length);
	const startRow = $derived(currentPage * pageSize + 1);
	const endRow = $derived(Math.min((currentPage + 1) * pageSize, totalRows));

	// Selection helpers
	const allSelected = $derived(paginatedRows.length > 0 && paginatedRows.every((_, i) => selectedRows.has(currentPage * pageSize + i)));
	const someSelected = $derived(paginatedRows.some((_, i) => selectedRows.has(currentPage * pageSize + i)));

	function handleSort(key: string, sortable: boolean = true) {
		if (!sortable) return;

		if (sortKey === key) {
			sortDirection = sortDirection === 'asc' ? 'desc' : 'asc';
		} else {
			sortKey = key;
			sortDirection = 'asc';
		}
		currentPage = 0;
	}

	function goToPage(page: number) {
		currentPage = Math.max(0, Math.min(page, totalPages - 1));
	}

	function formatCell(value: any, column: TableColumn): string {
		if (column.format) {
			return column.format(value);
		}
		if (value === null || value === undefined) {
			return '—';
		}
		if (typeof value === 'number') {
			// Auto-format numbers
			if (Math.abs(value) >= 1000000) {
				return (value / 1000000).toFixed(1) + 'M';
			}
			if (Math.abs(value) >= 1000) {
				return (value / 1000).toFixed(1) + 'K';
			}
			if (!Number.isInteger(value)) {
				return value.toFixed(2);
			}
		}
		if (value instanceof Date) {
			return value.toLocaleDateString();
		}
		return String(value);
	}

	function getCellClass(value: any, column: TableColumn): string {
		if (typeof value === 'number') {
			if (value > 0) return 'text-green-400';
			if (value < 0) return 'text-red-400';
		}
		return '';
	}

	function toggleColumnVisibility(key: string) {
		const newSet = new Set(visibleColumns);
		if (newSet.has(key)) {
			if (newSet.size > 1) {
				newSet.delete(key);
			}
		} else {
			newSet.add(key);
		}
		visibleColumns = newSet;
	}

	function toggleRowSelection(index: number) {
		const globalIndex = currentPage * pageSize + index;
		const newSet = new Set(selectedRows);
		if (newSet.has(globalIndex)) {
			newSet.delete(globalIndex);
		} else {
			newSet.add(globalIndex);
		}
		selectedRows = newSet;

		if (onSelectionChange) {
			const selected = Array.from(newSet).map((i) => sortedRows[i]).filter(Boolean);
			onSelectionChange(selected);
		}
	}

	function toggleSelectAll() {
		const newSet = new Set(selectedRows);
		const pageIndices = paginatedRows.map((_, i) => currentPage * pageSize + i);

		if (allSelected) {
			pageIndices.forEach((i) => newSet.delete(i));
		} else {
			pageIndices.forEach((i) => newSet.add(i));
		}
		selectedRows = newSet;

		if (onSelectionChange) {
			const selected = Array.from(newSet).map((i) => sortedRows[i]).filter(Boolean);
			onSelectionChange(selected);
		}
	}

	function handleRowClick(row: TableRow, index: number) {
		if (onRowClick) {
			onRowClick(row, currentPage * pageSize + index);
		}
	}

	// Column resizing
	function startResize(e: MouseEvent, columnKey: string) {
		e.preventDefault();
		e.stopPropagation();
		resizingColumn = columnKey;
		resizeStartX = e.clientX;
		resizeStartWidth = columnWidths[columnKey] || 150;

		window.addEventListener('mousemove', handleResize);
		window.addEventListener('mouseup', stopResize);
	}

	function handleResize(e: MouseEvent) {
		if (!resizingColumn) return;
		const diff = e.clientX - resizeStartX;
		const newWidth = Math.max(80, resizeStartWidth + diff);
		columnWidths = { ...columnWidths, [resizingColumn]: newWidth };
	}

	function stopResize() {
		resizingColumn = null;
		window.removeEventListener('mousemove', handleResize);
		window.removeEventListener('mouseup', stopResize);
	}

	// Export functions
	function exportCSV() {
		const headers = columns.map((c) => c.label).join(',');
		const csvRows = sortedRows.map((row) =>
			columns.map((col) => {
				const val = row[col.key];
				if (val === null || val === undefined) return '';
				const str = String(val);
				// Escape quotes and wrap in quotes if contains comma or quote
				if (str.includes(',') || str.includes('"') || str.includes('\n')) {
					return `"${str.replace(/"/g, '""')}"`;
				}
				return str;
			}).join(',')
		);

		const csv = [headers, ...csvRows].join('\n');
		downloadFile(csv, 'table-export.csv', 'text/csv');
	}

	function exportJSON() {
		const data = sortedRows.map((row) => {
			const obj: Record<string, any> = {};
			columns.forEach((col) => {
				obj[col.key] = row[col.key];
			});
			return obj;
		});
		downloadFile(JSON.stringify(data, null, 2), 'table-export.json', 'application/json');
	}

	function downloadFile(content: string, filename: string, mimeType: string) {
		const blob = new Blob([content], { type: mimeType });
		const url = URL.createObjectURL(blob);
		const a = document.createElement('a');
		a.href = url;
		a.download = filename;
		document.body.appendChild(a);
		a.click();
		document.body.removeChild(a);
		URL.revokeObjectURL(url);
	}

	function clearFilters() {
		searchQuery = '';
		columnFilters = {};
		currentPage = 0;
	}

	const hasActiveFilters = $derived(searchQuery.trim() !== '' || Object.values(columnFilters).some((v) => v.trim() !== ''));

	// Page size options
	const pageSizeOptions = [10, 25, 50, 100];

	// Pagination range
	const paginationRange = $derived.by(() => {
		const range: (number | 'ellipsis')[] = [];
		const delta = 2;

		for (let i = 0; i < totalPages; i++) {
			if (
				i === 0 ||
				i === totalPages - 1 ||
				(i >= currentPage - delta && i <= currentPage + delta)
			) {
				range.push(i);
			} else if (range[range.length - 1] !== 'ellipsis') {
				range.push('ellipsis');
			}
		}

		return range;
	});
</script>

<div
	bind:this={containerRef}
	class="table-chart-container rounded-lg border border-slate-700/50 bg-slate-800/20 overflow-hidden"
	class:compact={compactMode}
>
	<!-- Toolbar -->
	{#if searchable || exportable || showColumnToggle}
		<div class="table-toolbar flex items-center gap-3 p-3 border-b border-slate-700/50 bg-slate-800/40">
			<!-- Search -->
			{#if searchable}
				<div class="relative flex-1 max-w-xs">
					<Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
					<input
						type="text"
						bind:value={searchQuery}
						placeholder="Search all columns..."
						class="w-full pl-9 pr-8 py-1.5 rounded-lg bg-slate-900/50 border border-slate-600/50 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500"
					/>
					{#if searchQuery}
						<button
							class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 hover:text-white"
							onclick={() => (searchQuery = '')}
						>
							<X class="w-4 h-4" />
						</button>
					{/if}
				</div>
			{/if}

			<!-- Results count -->
			<div class="text-xs text-slate-400">
				{totalRows} {totalRows === 1 ? 'row' : 'rows'}
				{#if hasActiveFilters}
					<span class="text-indigo-400">(filtered)</span>
				{/if}
			</div>

			<div class="flex-1"></div>

			<!-- Clear filters -->
			{#if hasActiveFilters}
				<button
					class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs text-slate-300 hover:bg-slate-700/50 transition-colors"
					onclick={clearFilters}
				>
					<X class="w-3.5 h-3.5" />
					Clear filters
				</button>
			{/if}

			<!-- Column toggle -->
			{#if showColumnToggle}
				<div class="relative">
					<button
						class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs text-slate-300 hover:bg-slate-700/50 transition-colors border border-slate-600/50"
						onclick={() => (showColumnMenu = !showColumnMenu)}
					>
						<Settings2 class="w-3.5 h-3.5" />
						Columns
					</button>
					{#if showColumnMenu}
						<div class="absolute right-0 top-full mt-1 z-20 w-48 rounded-lg bg-slate-800 border border-slate-700 shadow-xl overflow-hidden">
							<div class="p-2 border-b border-slate-700 text-xs font-medium text-slate-400">
								Toggle Columns
							</div>
							<div class="max-h-64 overflow-y-auto p-1">
								{#each initialColumns as column}
									<button
										class="w-full flex items-center gap-2 px-2 py-1.5 rounded text-left text-sm text-slate-300 hover:bg-slate-700/50 transition-colors"
										onclick={() => toggleColumnVisibility(column.key)}
									>
										<div
											class="w-4 h-4 rounded border flex items-center justify-center {visibleColumns.has(column.key) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-500'}"
										>
											{#if visibleColumns.has(column.key)}
												<Check class="w-3 h-3 text-white" />
											{/if}
										</div>
										{column.label}
									</button>
								{/each}
							</div>
						</div>
					{/if}
				</div>
			{/if}

			<!-- Export -->
			{#if exportable}
				<div class="relative group">
					<button class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-lg text-xs text-slate-300 hover:bg-slate-700/50 transition-colors border border-slate-600/50">
						<Download class="w-3.5 h-3.5" />
						Export
					</button>
					<div class="absolute right-0 top-full mt-1 z-20 w-32 rounded-lg bg-slate-800 border border-slate-700 shadow-xl overflow-hidden opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all">
						<button
							class="w-full px-3 py-2 text-left text-sm text-slate-300 hover:bg-slate-700/50 transition-colors"
							onclick={exportCSV}
						>
							Export CSV
						</button>
						<button
							class="w-full px-3 py-2 text-left text-sm text-slate-300 hover:bg-slate-700/50 transition-colors"
							onclick={exportJSON}
						>
							Export JSON
						</button>
					</div>
				</div>
			{/if}
		</div>
	{/if}

	<!-- Table -->
	<div class="table-scroll-container overflow-auto" style={height ? `max-height: ${height}px` : ''}>
		<table bind:this={tableRef} class="w-full border-collapse">
			<thead class={stickyHeader ? 'sticky top-0 z-10' : ''}>
				<tr class="bg-slate-800/80 backdrop-blur-sm">
					<!-- Selection checkbox -->
					{#if selectable}
						<th class="table-header-cell w-10 px-3">
							<input
								type="checkbox"
								checked={allSelected}
								indeterminate={someSelected && !allSelected}
								onchange={toggleSelectAll}
								class="w-4 h-4 rounded border-slate-500 bg-slate-700 text-indigo-600 focus:ring-indigo-500 focus:ring-offset-0"
							/>
						</th>
					{/if}

					{#each columns as column}
						{@const width = columnWidths[column.key]}
						<th
							class="table-header-cell group relative select-none {column.sortable !== false ? 'cursor-pointer' : ''}"
							style="text-align: {column.align || 'left'}; {width ? `width: ${width}px; min-width: ${width}px;` : ''}"
							onclick={() => handleSort(column.key, column.sortable !== false)}
						>
							<div class="flex items-center gap-2 {column.align === 'center' ? 'justify-center' : column.align === 'right' ? 'justify-end' : 'justify-start'}">
								<span class="truncate">{column.label}</span>

								<!-- Sort indicator -->
								{#if column.sortable !== false}
									<span class="flex-shrink-0">
										{#if sortKey === column.key}
											{#if sortDirection === 'asc'}
												<ChevronUp class="w-4 h-4 text-indigo-400" />
											{:else}
												<ChevronDown class="w-4 h-4 text-indigo-400" />
											{/if}
										{:else}
											<ChevronsUpDown class="w-4 h-4 opacity-0 group-hover:opacity-50 transition-opacity" />
										{/if}
									</span>
								{/if}

								<!-- Column filter -->
								<button
									class="flex-shrink-0 p-0.5 rounded opacity-0 group-hover:opacity-100 hover:bg-slate-700/50 transition-all {columnFilters[column.key] ? 'opacity-100 text-indigo-400' : 'text-slate-400'}"
									onclick={(e) => {
										e.stopPropagation();
										showFilterMenu = showFilterMenu === column.key ? null : column.key;
									}}
								>
									<Filter class="w-3.5 h-3.5" />
								</button>
							</div>

							<!-- Column filter dropdown -->
							{#if showFilterMenu === column.key}
								<div
									class="absolute left-0 top-full mt-1 z-30 w-48 p-2 rounded-lg bg-slate-800 border border-slate-700 shadow-xl"
									onclick={(e) => e.stopPropagation()}
									onkeydown={(e) => e.stopPropagation()}
									role="dialog"
									aria-label="Filter options"
									tabindex="-1"
								>
									<input
										type="text"
										bind:value={columnFilters[column.key]}
										placeholder="Filter {column.label}..."
										class="w-full px-2 py-1.5 rounded bg-slate-900/50 border border-slate-600/50 text-sm text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500"
									/>
									{#if columnFilters[column.key]}
										<button
											class="mt-2 w-full px-2 py-1 rounded text-xs text-slate-400 hover:bg-slate-700/50 transition-colors"
											onclick={() => {
												columnFilters = { ...columnFilters, [column.key]: '' };
											}}
										>
											Clear filter
										</button>
									{/if}
								</div>
							{/if}

							<!-- Resize handle -->
							{#if resizable}
								<button
									type="button"
									class="absolute right-0 top-0 bottom-0 w-1 cursor-col-resize bg-transparent hover:bg-indigo-500/50 transition-colors border-0 p-0"
									onmousedown={(e) => startResize(e, column.key)}
									aria-label="Resize column"
								></button>
							{/if}
						</th>
					{/each}
				</tr>
			</thead>

			<tbody class="divide-y divide-slate-700/30">
				{#each paginatedRows as row, rowIndex}
					{@const globalIndex = currentPage * pageSize + rowIndex}
					{@const isSelected = selectedRows.has(globalIndex)}
					<tr
						class="table-row transition-colors {striped && rowIndex % 2 === 1 ? 'bg-slate-800/10' : ''} {hoverable ? 'hover:bg-slate-700/20' : ''} {isSelected ? 'bg-indigo-900/20' : ''} {onRowClick ? 'cursor-pointer' : ''}"
						onclick={() => handleRowClick(row, rowIndex)}
					>
						<!-- Selection checkbox -->
						{#if selectable}
							<td class="table-cell w-10 px-3" onclick={(e) => e.stopPropagation()}>
								<input
									type="checkbox"
									checked={isSelected}
									onchange={() => toggleRowSelection(rowIndex)}
									class="w-4 h-4 rounded border-slate-500 bg-slate-700 text-indigo-600 focus:ring-indigo-500 focus:ring-offset-0"
								/>
							</td>
						{/if}

						{#each columns as column}
							{@const value = row[column.key]}
							{@const width = columnWidths[column.key]}
							<td
								class="table-cell {getCellClass(value, column)}"
								style="text-align: {column.align || 'left'}; {width ? `width: ${width}px; min-width: ${width}px;` : ''}"
							>
								<span class="truncate block">{formatCell(value, column)}</span>
							</td>
						{/each}
					</tr>
				{/each}

				{#if paginatedRows.length === 0}
					<tr>
						<td colspan={columns.length + (selectable ? 1 : 0)} class="table-empty">
							<div class="flex flex-col items-center justify-center py-12">
								<svg class="w-12 h-12 text-slate-600 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
								</svg>
								<p class="text-slate-400 font-medium mb-1">No data found</p>
								{#if hasActiveFilters}
									<p class="text-slate-500 text-sm">Try adjusting your search or filters</p>
								{/if}
							</div>
						</td>
					</tr>
				{/if}
			</tbody>
		</table>
	</div>

	<!-- Footer with pagination -->
	{#if paginate && rows.length > 0}
		<div class="table-footer flex items-center justify-between px-4 py-3 border-t border-slate-700/50 bg-slate-800/40">
			<!-- Page size selector -->
			<div class="flex items-center gap-2">
				<span class="text-xs text-slate-400">Show</span>
				<select
					bind:value={pageSize}
					onchange={() => (currentPage = 0)}
					class="px-2 py-1 rounded bg-slate-900/50 border border-slate-600/50 text-xs text-white focus:outline-none focus:border-indigo-500"
				>
					{#each pageSizeOptions as option}
						<option value={option}>{option}</option>
					{/each}
				</select>
				<span class="text-xs text-slate-400">rows</span>
			</div>

			<!-- Results info -->
			<div class="text-xs text-slate-400">
				{#if totalRows > 0}
					{startRow}-{endRow} of {totalRows}
				{:else}
					No results
				{/if}
			</div>

			<!-- Pagination controls -->
			<div class="flex items-center gap-1">
				<button
					class="p-1.5 rounded hover:bg-slate-700/50 text-slate-400 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
					disabled={currentPage === 0}
					onclick={() => goToPage(0)}
					title="First page"
				>
					<ChevronsLeft class="w-4 h-4" />
				</button>
				<button
					class="p-1.5 rounded hover:bg-slate-700/50 text-slate-400 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
					disabled={currentPage === 0}
					onclick={() => goToPage(currentPage - 1)}
					title="Previous page"
				>
					<ChevronLeft class="w-4 h-4" />
				</button>

				{#each paginationRange as page}
					{#if page === 'ellipsis'}
						<span class="px-2 text-slate-500">...</span>
					{:else}
						<button
							class="min-w-[28px] h-7 px-2 rounded text-xs font-medium transition-colors {currentPage === page ? 'bg-indigo-600 text-white' : 'text-slate-400 hover:bg-slate-700/50 hover:text-white'}"
							onclick={() => goToPage(page)}
						>
							{page + 1}
						</button>
					{/if}
				{/each}

				<button
					class="p-1.5 rounded hover:bg-slate-700/50 text-slate-400 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
					disabled={currentPage >= totalPages - 1}
					onclick={() => goToPage(currentPage + 1)}
					title="Next page"
				>
					<ChevronRight class="w-4 h-4" />
				</button>
				<button
					class="p-1.5 rounded hover:bg-slate-700/50 text-slate-400 hover:text-white disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
					disabled={currentPage >= totalPages - 1}
					onclick={() => goToPage(totalPages - 1)}
					title="Last page"
				>
					<ChevronsRight class="w-4 h-4" />
				</button>
			</div>
		</div>
	{/if}
</div>

<!-- Click outside handlers -->
<svelte:window
	onclick={() => {
		showColumnMenu = false;
		showFilterMenu = null;
	}}
/>
