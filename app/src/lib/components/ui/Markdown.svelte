<script lang="ts">
	import { marked } from 'marked';

	interface Props {
		content: string;
		class?: string;
	}

	let { content, class: className = '' }: Props = $props();

	// Configure marked for safe rendering
	marked.setOptions({
		breaks: true,
		gfm: true
	});

	let html = $derived(marked.parse(content || '') as string);
</script>

<div class="markdown-content {className}">
	{@html html}
</div>

<style>
@reference "$src/app.css";

@layer components.markdown {
	.markdown-content {
		@apply text-base-200 text-sm leading-relaxed;
	}

	.markdown-content :global(h1) {
		@apply text-2xl font-bold text-white mt-6 mb-4 first:mt-0;
	}

	.markdown-content :global(h2) {
		@apply text-xl font-semibold text-white mt-6 mb-3 first:mt-0;
	}

	.markdown-content :global(h3) {
		@apply text-lg font-semibold text-white mt-5 mb-2 first:mt-0;
	}

	.markdown-content :global(h4) {
		@apply text-base font-semibold text-white mt-4 mb-2 first:mt-0;
	}

	.markdown-content :global(p) {
		@apply mb-4 last:mb-0;
	}

	.markdown-content :global(ul) {
		@apply list-disc list-inside mb-4 space-y-1;
	}

	.markdown-content :global(ol) {
		@apply list-decimal list-inside mb-4 space-y-1;
	}

	.markdown-content :global(li) {
		@apply text-base-300;
	}

	.markdown-content :global(li > ul),
	.markdown-content :global(li > ol) {
		@apply ml-4 mt-1;
	}

	.markdown-content :global(strong) {
		@apply font-semibold text-white;
	}

	.markdown-content :global(em) {
		@apply italic;
	}

	.markdown-content :global(code) {
		@apply px-1.5 py-0.5 bg-base-700 text-primary rounded text-xs font-mono;
	}

	.markdown-content :global(pre) {
		@apply p-4 bg-base-800 rounded-lg overflow-x-auto mb-4;
	}

	.markdown-content :global(pre code) {
		@apply p-0 bg-transparent rounded-none text-base-200;
	}

	.markdown-content :global(blockquote) {
		@apply border-l-4 border-primary/50 pl-4 italic text-base-400 my-4;
	}

	.markdown-content :global(a) {
		@apply text-primary hover:text-primary-light underline;
	}

	.markdown-content :global(hr) {
		@apply border-base-700 my-6;
	}

	.markdown-content :global(table) {
		@apply w-full border-collapse mb-4;
	}

	.markdown-content :global(th) {
		@apply text-left p-2 bg-base-800 border border-base-700 text-white font-medium;
	}

	.markdown-content :global(td) {
		@apply p-2 border border-base-700;
	}

	.markdown-content :global(img) {
		@apply max-w-full h-auto rounded-lg my-4;
	}
}
</style>
