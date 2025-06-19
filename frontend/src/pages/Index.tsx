
const Index = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-50 via-blue-50 to-indigo-100 animate-fade-in">
      <div className="text-center max-w-2xl mx-auto px-6 animate-scale-in">
        <div className="mb-8 transform transition-all duration-500 hover:scale-110">
          <div className="w-32 h-32 bg-gradient-to-r from-blue-600 to-purple-600 rounded-full mx-auto mb-6 flex items-center justify-center shadow-2xl hover:shadow-3xl transition-all duration-500 hover:rotate-12">
            <div className="w-16 h-16 bg-white rounded-full flex items-center justify-center">
              <span className="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                A
              </span>
            </div>
          </div>
        </div>
        
        <h1 className="text-6xl font-bold mb-6 bg-gradient-to-r from-blue-600 via-purple-600 to-indigo-600 bg-clip-text text-transparent animate-fade-in hover:scale-105 transition-transform duration-300">
          Welcome to Your Blank App
        </h1>
        
        <p className="text-xl text-gray-600 mb-8 animate-fade-in leading-relaxed" style={{ animationDelay: '200ms' }}>
          Start building your amazing project here! The possibilities are endless.
        </p>
        
        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center animate-fade-in" style={{ animationDelay: '400ms' }}>
          <button className="px-8 py-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-full font-semibold shadow-lg hover:shadow-2xl transition-all duration-300 hover:scale-105 active:scale-95 hover:-translate-y-1 relative overflow-hidden group">
            <span className="relative z-10">Get Started</span>
            <div className="absolute inset-0 bg-gradient-to-r from-purple-600 to-blue-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          </button>
          
          <button className="px-8 py-4 border-2 border-gray-300 text-gray-700 rounded-full font-semibold hover:border-blue-500 hover:text-blue-600 transition-all duration-300 hover:scale-105 active:scale-95 hover:shadow-lg relative overflow-hidden group">
            <span className="relative z-10">Learn More</span>
            <div className="absolute inset-0 bg-blue-50 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          </button>
        </div>
        
        <div className="mt-16 grid grid-cols-1 md:grid-cols-3 gap-8 animate-fade-in" style={{ animationDelay: '600ms' }}>
          <div className="p-6 bg-white rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 transform hover:scale-105 group">
            <div className="w-12 h-12 bg-gradient-to-r from-blue-500 to-cyan-500 rounded-lg mx-auto mb-4 flex items-center justify-center group-hover:rotate-12 transition-transform duration-300">
              <span className="text-white font-bold text-xl">⚡</span>
            </div>
            <h3 className="text-xl font-semibold mb-2 text-gray-800 group-hover:text-blue-600 transition-colors duration-300">Fast</h3>
            <p className="text-gray-600">Lightning-fast development and deployment</p>
          </div>
          
          <div className="p-6 bg-white rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 transform hover:scale-105 group">
            <div className="w-12 h-12 bg-gradient-to-r from-purple-500 to-pink-500 rounded-lg mx-auto mb-4 flex items-center justify-center group-hover:rotate-12 transition-transform duration-300">
              <span className="text-white font-bold text-xl">🎨</span>
            </div>
            <h3 className="text-xl font-semibold mb-2 text-gray-800 group-hover:text-purple-600 transition-colors duration-300">Beautiful</h3>
            <p className="text-gray-600">Stunning designs that captivate users</p>
          </div>
          
          <div className="p-6 bg-white rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-500 hover:-translate-y-2 transform hover:scale-105 group">
            <div className="w-12 h-12 bg-gradient-to-r from-green-500 to-emerald-500 rounded-lg mx-auto mb-4 flex items-center justify-center group-hover:rotate-12 transition-transform duration-300">
              <span className="text-white font-bold text-xl">🚀</span>
            </div>
            <h3 className="text-xl font-semibold mb-2 text-gray-800 group-hover:text-green-600 transition-colors duration-300">Powerful</h3>
            <p className="text-gray-600">Robust features for any application</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Index;
