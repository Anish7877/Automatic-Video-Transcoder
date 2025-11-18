import { useState } from 'react';
import { Upload, Video, FileVideo, Loader2 } from 'lucide-react';
import FileUpload from './components/FileUpload';
import FormatSelector from './components/FormatSelector';
import TranscodingStatus from './components/TranscodingStatus';

function App() {
  const [file, setFile] = useState<File | null>(null);
  const [selectedFormat, setSelectedFormat] = useState<string>('mp4');
  const [isProcessing, setIsProcessing] = useState(false);

  const handleFileSelect = (selectedFile: File | null) => {
    setFile(selectedFile);
    setIsProcessing(false);
  };

  

const handleTranscode = async () => {
  if (!file) return;
  setIsProcessing(true);

  try {
    const formData = new FormData();
    formData.append('video', file);
    formData.append('format', selectedFormat);

    // Replace with your actual backend URL
    const response = await fetch('http://localhost:8080/upload', {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Conversion failed');
    }

    const data = await response.json();
    
    // Handle success - typically you might want to set a download URL state here
    console.log("Download URL:", data.downloadUrl);
    
    // Automatically trigger download if desired
    window.location.href = data.downloadUrl;

    } catch (error) {
      console.error("Error transcoding video:", error);
      alert("An error occurred during transcoding.");
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100">
      <div className="container mx-auto px-4 py-12">
        <header className="text-center mb-16">
          <div className="flex items-center justify-center mb-4">
            <FileVideo className="w-12 h-12 text-slate-800" />
          </div>
          <h1 className="text-4xl font-light text-slate-900 mb-3 tracking-tight">
            Video Transcoder
          </h1>
          <p className="text-slate-600 text-lg font-light">
            Convert your videos to any format with ease
          </p>
        </header>

        <div className="max-w-3xl mx-auto">
          <div className="bg-white rounded-2xl shadow-sm border border-slate-200 overflow-hidden">
            <div className="p-8">
              <FileUpload
                onFileSelect={handleFileSelect}
                currentFile={file}
                disabled={isProcessing}
              />

              {file && !isProcessing && (
                <div className="mt-8 space-y-6">
                  <FormatSelector
                    selectedFormat={selectedFormat}
                    onFormatChange={setSelectedFormat}
                  />

                  <button
                    onClick={handleTranscode}
                    className="w-full bg-slate-900 hover:bg-slate-800 text-white py-4 rounded-xl font-medium transition-colors duration-200 flex items-center justify-center gap-2"
                  >
                    <Video className="w-5 h-5" />
                    Convert to {selectedFormat.toUpperCase()}
                  </button>
                </div>
              )}

              {isProcessing && (
                <TranscodingStatus
                  fileName={file?.name || ''}
                  targetFormat={selectedFormat}
                />
              )}
            </div>
          </div>

          <div className="mt-8 text-center text-sm text-slate-500 font-light">
            Supports MP4, AVI, MOV, MKV, WebM, and more
          </div>
        </div>
      </div>
    </div>
  );
}

export default App;
