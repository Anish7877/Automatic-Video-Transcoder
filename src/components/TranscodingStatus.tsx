// import { Loader2, CheckCircle } from 'lucide-react';
// import { useState, useEffect } from 'react';

// interface TranscodingStatusProps {
//   fileName: string;
//   targetFormat: string;
// }

// function TranscodingStatus({ fileName, targetFormat }: TranscodingStatusProps) {
//   const [progress, setProgress] = useState(0);
//   const [isComplete, setIsComplete] = useState(false);

//   useEffect(() => {
//     const interval = setInterval(() => {
//       setProgress((prev) => {
//         if (prev >= 100) {
//           clearInterval(interval);
//           setIsComplete(true);
//           return 100;
//         }
//         return prev + 2;
//       });
//     }, 100);

//     return () => clearInterval(interval);
//   }, []);

//   return (
//     <div className="mt-8">
//       <div className="flex items-center justify-between mb-4">
//         <div className="flex items-center gap-3">
//           {!isComplete ? (
//             <Loader2 className="w-6 h-6 text-slate-600 animate-spin" />
//           ) : (
//             <CheckCircle className="w-6 h-6 text-green-600" />
//           )}
//           <div>
//             <p className="font-medium text-slate-900">
//               {isComplete ? 'Conversion complete' : 'Converting video...'}
//             </p>
//             <p className="text-sm text-slate-500 font-light">
//               {fileName} → {targetFormat.toUpperCase()}
//             </p>
//           </div>
//         </div>
//         <span className="text-sm font-medium text-slate-600">
//           {progress}%
//         </span>
//       </div>

//       <div className="w-full bg-slate-200 rounded-full h-2 overflow-hidden">
//         <div
//           className="bg-slate-900 h-full transition-all duration-300 ease-out rounded-full"
//           style={{ width: `${progress}%` }}
//         />
//       </div>

//       {isComplete && (
//         <button className="mt-6 w-full bg-green-600 hover:bg-green-700 text-white py-4 rounded-xl font-medium transition-colors duration-200">
//           Download Converted File
//         </button>
//       )}
//     </div>
//   );
// }

// export default TranscodingStatus;

import { Loader2, CheckCircle } from 'lucide-react';
import { useState, useEffect } from 'react';

interface TranscodingStatusProps {
  fileName: string;
  targetFormat: string;
}

function TranscodingStatus({ fileName, targetFormat }: TranscodingStatusProps) {
  const [progress, setProgress] = useState(0);
  // isComplete is intentionally designed to never become true based on this timer
  const [isComplete, setIsComplete] = useState(false); 

  useEffect(() => {
    const interval = setInterval(() => {
      setProgress((prev) => {
        // Stop the animation when the progress hits 98 (so it returns 99)
        if (prev >= 98) {
          clearInterval(interval);
          // *** LOOPHOLE IMPLEMENTATION: We stop the interval but DO NOT set setIsComplete(true) ***
          return 99; // Lock the visible progress at 99%
        }
        return prev + 2;
      });
    }, 100);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="mt-8">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          {/* This will always show the Loader2 because isComplete is false */}
          <Loader2 className="w-6 h-6 text-slate-600 animate-spin" />
          
          <div>
            <p className="font-medium text-slate-900">
              {/* This status text will always show 'Converting video...' */}
              Converting video...
            </p>
            <p className="text-sm text-slate-500 font-light">
              {fileName} → {targetFormat.toUpperCase()}
            </p>
          </div>
        </div>
        <span className="text-sm font-medium text-slate-600">
          {progress}%
        </span>
      </div>

      <div className="w-full bg-slate-200 rounded-full h-2 overflow-hidden">
        <div
          className="bg-slate-900 h-full transition-all duration-300 ease-out rounded-full"
          style={{ width: `${progress}%` }}
        />
      </div>

      {/* This button will NEVER appear because isComplete will remain false */}
      {isComplete && (
        <button className="mt-6 w-full bg-green-600 hover:bg-green-700 text-white py-4 rounded-xl font-medium transition-colors duration-200">
          Download Converted File
        </button>
      )}
    </div>
  );
}

export default TranscodingStatus;