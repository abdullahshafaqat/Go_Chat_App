
import { useLocation } from "react-router-dom";
import { useEffect } from "react";

const NotFound = () => {
  const location = useLocation();

  useEffect(() => {
    console.error(
      "404 Error: User attempted to access non-existent route:",
      location.pathname
    );
  }, [location.pathname]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-red-50 via-orange-50 to-yellow-50 animate-fade-in">
      <div className="text-center max-w-md mx-auto px-6 animate-scale-in">
        <div className="mb-8 transform transition-all duration-500 hover:scale-110">
          <div className="text-9xl font-bold text-transparent bg-gradient-to-r from-red-500 to-orange-500 bg-clip-text animate-pulse hover:animate-bounce">
            404
          </div>
        </div>
        
        <div className="mb-6 animate-fade-in" style={{ animationDelay: '200ms' }}>
          <h1 className="text-3xl font-bold mb-2 text-gray-800 hover:text-red-600 transition-colors duration-300">
            Oops! Page not found
          </h1>
          <p className="text-lg text-gray-600 leading-relaxed">
            The page you're looking for seems to have wandered off into the digital void.
          </p>
        </div>
        
        <div className="animate-fade-in" style={{ animationDelay: '400ms' }}>
          <a 
            href="/" 
            className="inline-flex items-center px-6 py-3 bg-gradient-to-r from-red-500 to-orange-500 text-white rounded-full font-semibold shadow-lg hover:shadow-2xl transition-all duration-300 hover:scale-105 active:scale-95 hover:-translate-y-1 relative overflow-hidden group"
          >
            <span className="relative z-10 flex items-center">
              <svg className="w-5 h-5 mr-2 group-hover:-translate-x-1 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
              </svg>
              Return to Home
            </span>
            <div className="absolute inset-0 bg-gradient-to-r from-orange-500 to-red-500 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          </a>
        </div>
        
        <div className="mt-12 animate-fade-in" style={{ animationDelay: '600ms' }}>
          <div className="grid grid-cols-2 gap-4 text-center">
            <div className="p-4 bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-300 hover:scale-105 transform group">
              <div className="text-2xl mb-2 group-hover:animate-bounce">🏠</div>
              <p className="text-sm text-gray-600">Go Home</p>
            </div>
            <div className="p-4 bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-300 hover:scale-105 transform group">
              <div className="text-2xl mb-2 group-hover:animate-spin">🔍</div>
              <p className="text-sm text-gray-600">Search</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default NotFound;
