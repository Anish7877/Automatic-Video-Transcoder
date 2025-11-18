import { useRef } from 'react';
import { Upload, X, FileVideo } from 'lucide-react';

interface FileUploadProps {
  onFileSelect: (file: File | null) => void;
  currentFile: File | null;
  disabled?: boolean;
}

function FileUpload({ onFileSelect, currentFile, disabled }: FileUploadProps) {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = event.target.files?.[0];
    if (selectedFile && selectedFile.type.startsWith('video/')) {
      onFileSelect(selectedFile);
    }
  };

  const handleDrop = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
    const droppedFile = event.dataTransfer.files[0];
    if (droppedFile && droppedFile.type.startsWith('video/')) {
      onFileSelect(droppedFile);
    }
  };

  const handleDragOver = (event: React.DragEvent<HTMLDivElement>) => {
    event.preventDefault();
  };

  const handleRemoveFile = () => {
    onFileSelect(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  return (
    <div>
      {!currentFile ? (
        <div
          onDrop={handleDrop}
          onDragOver={handleDragOver}
          className={`border-2 border-dashed border-slate-300 rounded-xl p-12 text-center transition-all duration-200 ${
            disabled
              ? 'opacity-50 cursor-not-allowed'
              : 'hover:border-slate-400 hover:bg-slate-50 cursor-pointer'
          }`}
          onClick={() => !disabled && fileInputRef.current?.click()}
        >
          <div className="flex flex-col items-center gap-4">
            <div className="w-16 h-16 rounded-full bg-slate-100 flex items-center justify-center">
              <Upload className="w-8 h-8 text-slate-600" />
            </div>
            <div>
              <p className="text-lg font-medium text-slate-900 mb-1">
                Drop your video here
              </p>
              <p className="text-sm text-slate-500 font-light">
                or click to browse files
              </p>
            </div>
            <p className="text-xs text-slate-400 font-light">
              Supports MP4, AVI, MOV, MKV, WebM and other video formats
            </p>
          </div>
          <input
            ref={fileInputRef}
            type="file"
            accept="video/*"
            onChange={handleFileChange}
            className="hidden"
            disabled={disabled}
          />
        </div>
      ) : (
        <div className="border border-slate-200 rounded-xl p-6 bg-slate-50">
          <div className="flex items-start justify-between">
            <div className="flex items-start gap-4 flex-1">
              <div className="w-12 h-12 rounded-lg bg-slate-200 flex items-center justify-center flex-shrink-0">
                <FileVideo className="w-6 h-6 text-slate-600" />
              </div>
              <div className="flex-1 min-w-0">
                <p className="font-medium text-slate-900 truncate">
                  {currentFile.name}
                </p>
                <p className="text-sm text-slate-500 mt-1 font-light">
                  {formatFileSize(currentFile.size)}
                </p>
              </div>
            </div>
            {!disabled && (
              <button
                onClick={handleRemoveFile}
                className="ml-4 p-2 hover:bg-slate-200 rounded-lg transition-colors duration-200"
              >
                <X className="w-5 h-5 text-slate-600" />
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}

export default FileUpload;
