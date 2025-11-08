import { useState, useCallback } from "react";
import { VideoUpload } from "@/components/VideoUpload";
import { FormatSelector } from "@/components/FormatSelector";
import { ConversionProgress } from "@/components/ConversionProgress";
import { ThemeToggle } from "@/components/ThemeToggle";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { useToast } from "@/hooks/use-toast";
import { Download, Video } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

export default function Home() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [selectedFormat, setSelectedFormat] = useState("mp4");
  const [isConverting, setIsConverting] = useState(false);
  const [isComplete, setIsComplete] = useState(false);
  const [progress, setProgress] = useState(0);
  const [convertedFileUrl, setConvertedFileUrl] = useState<string | null>(null);
  const { toast } = useToast();

  const handleFileSelect = useCallback((file: File | null) => {
    setSelectedFile(file);
    setIsComplete(false);
    setProgress(0);
    setConvertedFileUrl(null);

    if (file) {
      toast({
        title: "Upload successful",
        description: `${file.name} has been uploaded.`,
      });
    }
  }, [toast]);

  const handleConvert = useCallback(() => {
    if (!selectedFile) return;

    setIsConverting(true);
    setIsComplete(false);
    setProgress(0);

    toast({
      title: "Conversion started",
      description: `Converting to ${selectedFormat.toUpperCase()}...`,
    });

    const interval = setInterval(() => {
      setProgress((prev) => {
        if (prev >= 100) {
          clearInterval(interval);
          setIsConverting(false);
          setIsComplete(true);

          const blob = new Blob([selectedFile], { type: `video/${selectedFormat}` });
          const url = URL.createObjectURL(blob);
          setConvertedFileUrl(url);

          toast({
            title: "Conversion complete",
            description: "Your video is ready to download!",
          });

          return 100;
        }
        return prev + 2;
      });
    }, 50);
  }, [selectedFile, selectedFormat, toast]);

  const handleDownload = useCallback(() => {
    if (!convertedFileUrl || !selectedFile) return;

    const link = document.createElement("a");
    link.href = convertedFileUrl;
    link.download = `${selectedFile.name.split(".")[0]}_converted.${selectedFormat}`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    toast({
      title: "Download started",
      description: "Your converted video is downloading.",
    });
  }, [convertedFileUrl, selectedFile, selectedFormat, toast]);

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <header className="border-b border-border bg-card/50 backdrop-blur-sm sticky top-0 z-50">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="bg-primary text-primary-foreground p-2 rounded-md">
              <Video className="h-6 w-6" />
            </div>
            <h1 className="text-xl sm:text-2xl font-bold">
              Automatic Video Transcoder
            </h1>
          </div>
          <ThemeToggle />
        </div>
      </header>

      <main className="flex-1 container mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 lg:py-16 flex items-center justify-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
          className="w-full max-w-3xl"
        >
          <Card className="shadow-lg">
            <CardContent className="p-6 sm:p-8 lg:p-10 space-y-8">
              <div className="space-y-2 text-center">
                <h2 className="text-2xl sm:text-3xl font-bold">
                  Convert Your Video
                </h2>
                <p className="text-muted-foreground">
                  Upload a video, select your desired format, and download the
                  converted file
                </p>
              </div>

              <VideoUpload
                selectedFile={selectedFile}
                onFileSelect={handleFileSelect}
              />

              <AnimatePresence>
                {selectedFile && !isComplete && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    exit={{ opacity: 0, height: 0 }}
                    transition={{ duration: 0.3 }}
                  >
                    <FormatSelector
                      selectedFormat={selectedFormat}
                      onFormatChange={setSelectedFormat}
                    />
                  </motion.div>
                )}
              </AnimatePresence>

              <ConversionProgress
                progress={progress}
                isConverting={isConverting}
                isComplete={isComplete}
              />

              <AnimatePresence>
                {selectedFile && !isConverting && !isComplete && (
                  <motion.div
                    initial={{ opacity: 0, height: 0 }}
                    animate={{ opacity: 1, height: "auto" }}
                    exit={{ opacity: 0, height: 0 }}
                    transition={{ duration: 0.3 }}
                  >
                    <Button
                      onClick={handleConvert}
                      className="w-full"
                      size="lg"
                      data-testid="button-convert"
                    >
                      Convert to {selectedFormat.toUpperCase()}
                    </Button>
                  </motion.div>
                )}

                {isComplete && (
                  <motion.div
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    transition={{ duration: 0.3 }}
                  >
                    <Button
                      onClick={handleDownload}
                      className="w-full"
                      size="lg"
                      data-testid="button-download"
                    >
                      <Download className="mr-2 h-5 w-5" />
                      Download Converted File
                    </Button>
                  </motion.div>
                )}
              </AnimatePresence>
            </CardContent>
          </Card>
        </motion.div>
      </main>
    </div>
  );
}
