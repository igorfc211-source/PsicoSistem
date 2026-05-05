export const MAX_DOCUMENT_BYTES = 20 * 1024 * 1024;

const DOCUMENT_DB_NAME = 'psicosistem-documents';
const DOCUMENT_STORE_NAME = 'files';

// Cria uma chave estavel por aprendente para evitar colisao entre documentos.
export function getDocumentStorageKey(learnerId: string, documentId: string) {
	return `${learnerId}:${documentId}`;
}

// Prepara o arquivo para armazenamento e tenta compactar automaticamente quando passa de 20 MB.
export async function prepareDocumentBlob(file: File) {
	if (file.size <= MAX_DOCUMENT_BYTES) {
		return {
			blob: file,
			compressed: false
		};
	}

	if (!('CompressionStream' in globalThis)) {
		throw new Error('Este navegador nao suporta compressao automatica de arquivos grandes.');
	}

	const compressedStream = file.stream().pipeThrough(new CompressionStream('gzip'));
	const compressedBlob = await new Response(compressedStream).blob();

	if (compressedBlob.size > MAX_DOCUMENT_BYTES) {
		throw new Error('Mesmo comprimido, o documento continua acima de 20 MB.');
	}

	return {
		blob: compressedBlob,
		compressed: true
	};
}

// Abre o banco local usado para documentos, isolando IndexedDB da camada de interface.
export function openDocumentDb() {
	return new Promise<IDBDatabase>((resolve, reject) => {
		const request = indexedDB.open(DOCUMENT_DB_NAME, 1);

		request.onupgradeneeded = () => {
			request.result.createObjectStore(DOCUMENT_STORE_NAME);
		};

		request.onsuccess = () => resolve(request.result);
		request.onerror = () => reject(request.error);
	});
}

// Salva o blob do documento sem inflar o localStorage com arquivos grandes.
export async function putDocumentBlob(id: string, blob: Blob) {
	const db = await openDocumentDb();

	return new Promise<void>((resolve, reject) => {
		const transaction = db.transaction(DOCUMENT_STORE_NAME, 'readwrite');
		transaction.objectStore(DOCUMENT_STORE_NAME).put(blob, id);
		transaction.oncomplete = () => resolve();
		transaction.onerror = () => reject(transaction.error);
	});
}

// Recupera o arquivo original/compactado para download pelo usuario.
export async function getDocumentBlob(id: string) {
	const db = await openDocumentDb();

	return new Promise<Blob | null>((resolve, reject) => {
		const transaction = db.transaction(DOCUMENT_STORE_NAME, 'readonly');
		const request = transaction.objectStore(DOCUMENT_STORE_NAME).get(id);
		request.onsuccess = () => resolve((request.result as Blob | undefined) ?? null);
		request.onerror = () => reject(request.error);
	});
}

// Remove o blob do IndexedDB quando o documento e excluido do prontuario.
export async function deleteDocumentBlob(id: string) {
	const db = await openDocumentDb();

	return new Promise<void>((resolve, reject) => {
		const transaction = db.transaction(DOCUMENT_STORE_NAME, 'readwrite');
		transaction.objectStore(DOCUMENT_STORE_NAME).delete(id);
		transaction.oncomplete = () => resolve();
		transaction.onerror = () => reject(transaction.error);
	});
}
