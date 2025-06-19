
import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useAuth } from '../contexts/AuthContext';

const Signup: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [passwordStrength, setPasswordStrength] = useState(0);
  const { signup, isLoading, isAuthenticated } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/chat', { replace: true });
    }
  }, [isAuthenticated, navigate]);

  const validateEmail = (email: string) => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  };

  const calculatePasswordStrength = (password: string) => {
    let strength = 0;
    if (password.length >= 8) strength += 1;
    if (/[A-Z]/.test(password)) strength += 1;
    if (/[a-z]/.test(password)) strength += 1;
    if (/\d/.test(password)) strength += 1;
    if (/[!@#$%^&*]/.test(password)) strength += 1;
    return strength;
  };

  const handleInputChange = (field: string, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    
    if (field === 'password') {
      setPasswordStrength(calculatePasswordStrength(value));
    }
  };

  const getPasswordStrengthColor = () => {
    if (passwordStrength <= 2) return 'bg-red-500';
    if (passwordStrength <= 3) return 'bg-yellow-500';
    return 'bg-green-500';
  };

  const getPasswordStrengthText = () => {
    if (passwordStrength <= 2) return 'Weak';
    if (passwordStrength <= 3) return 'Medium';
    return 'Strong';
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const newErrors: Record<string, string> = {};
    
    if (!formData.username.trim()) {
      newErrors.username = 'Username is required';
    }
    
    if (!formData.email) {
      newErrors.email = 'Email is required';
    } else if (!validateEmail(formData.email)) {
      newErrors.email = 'Please enter a valid email';
    }
    
    if (!formData.password) {
      newErrors.password = 'Password is required';
    } else if (formData.password.length < 6) {
      newErrors.password = 'Password must be at least 6 characters';
    }
    
    if (!formData.confirmPassword) {
      newErrors.confirmPassword = 'Please confirm your password';
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = 'Passwords do not match';
    }
    
    setErrors(newErrors);
    
    if (Object.keys(newErrors).length === 0) {
      try {
        await signup(formData.username, formData.email, formData.password);
        navigate('/chat');
      } catch (error) {
        // Error handling is done in the auth context
      }
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-50 via-pink-50 to-indigo-100 animate-fade-in relative overflow-hidden">
      {/* Animated background elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-4 -left-4 w-72 h-72 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-float"></div>
        <div className="absolute -top-4 -right-4 w-72 h-72 bg-pink-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-float" style={{ animationDelay: '2s' }}></div>
        <div className="absolute -bottom-8 left-20 w-72 h-72 bg-indigo-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-float" style={{ animationDelay: '4s' }}></div>
      </div>

      <div className="relative max-w-md w-full space-y-8 p-8 bg-white/80 backdrop-blur-lg rounded-2xl shadow-xl hover:shadow-2xl transition-all duration-500 transform hover:-translate-y-1 animate-scale-in border border-white/20">
        <div className="text-center animate-fade-in">
          <div className="mb-4 transform transition-all duration-500 hover:scale-110">
            <div className="w-16 h-16 bg-gradient-to-r from-purple-600 to-pink-600 rounded-full mx-auto flex items-center justify-center shadow-lg hover:shadow-xl transition-all duration-300 animate-pulse hover:animate-none">
              <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
          </div>
          <h2 className="text-3xl font-extrabold text-gray-900 animate-fade-in bg-gradient-to-r from-purple-600 to-pink-600 bg-clip-text text-transparent">
            Join Our Community
          </h2>
          <p className="mt-2 text-sm text-gray-600 animate-fade-in delay-150">
            Create your account to get started
          </p>
        </div>
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <div className="transform transition-all duration-300 hover:scale-105 animate-slide-in-right" style={{ animationDelay: '100ms' }}>
              <Label htmlFor="username" className="text-gray-700 font-medium">Username</Label>
              <Input
                id="username"
                type="text"
                value={formData.username}
                onChange={(e) => handleInputChange('username', e.target.value)}
                className={`mt-1 transition-all duration-300 focus:scale-105 hover:shadow-md ${errors.username ? 'border-red-500 animate-pulse' : 'focus:border-purple-500'}`}
                placeholder="Enter your username"
              />
              {errors.username && (
                <p className="mt-1 text-sm text-red-600 animate-fade-in">{errors.username}</p>
              )}
            </div>
            
            <div className="transform transition-all duration-300 hover:scale-105 animate-slide-in-right" style={{ animationDelay: '200ms' }}>
              <Label htmlFor="email" className="text-gray-700 font-medium">Email</Label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => handleInputChange('email', e.target.value)}
                className={`mt-1 transition-all duration-300 focus:scale-105 hover:shadow-md ${errors.email ? 'border-red-500 animate-pulse' : 'focus:border-purple-500'}`}
                placeholder="Enter your email"
              />
              {errors.email && (
                <p className="mt-1 text-sm text-red-600 animate-fade-in">{errors.email}</p>
              )}
            </div>
            
            <div className="transform transition-all duration-300 hover:scale-105 animate-slide-in-right" style={{ animationDelay: '300ms' }}>
              <Label htmlFor="password" className="text-gray-700 font-medium">Password</Label>
              <Input
                id="password"
                type="password"
                value={formData.password}
                onChange={(e) => handleInputChange('password', e.target.value)}
                className={`mt-1 transition-all duration-300 focus:scale-105 hover:shadow-md ${errors.password ? 'border-red-500 animate-pulse' : 'focus:border-purple-500'}`}
                placeholder="Enter your password"
              />
              {formData.password && (
                <div className="mt-2 animate-fade-in">
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-200 rounded-full h-2 overflow-hidden">
                      <div
                        className={`h-2 rounded-full transition-all duration-500 ${getPasswordStrengthColor()} animate-scale-in`}
                        style={{ width: `${(passwordStrength / 5) * 100}%` }}
                      ></div>
                    </div>
                    <span className={`text-xs font-medium transition-colors duration-300 ${getPasswordStrengthColor().replace('bg-', 'text-')}`}>
                      {getPasswordStrengthText()}
                    </span>
                  </div>
                </div>
              )}
              {errors.password && (
                <p className="mt-1 text-sm text-red-600 animate-fade-in">{errors.password}</p>
              )}
            </div>
            
            <div className="transform transition-all duration-300 hover:scale-105 animate-slide-in-right" style={{ animationDelay: '400ms' }}>
              <Label htmlFor="confirmPassword" className="text-gray-700 font-medium">Confirm Password</Label>
              <Input
                id="confirmPassword"
                type="password"
                value={formData.confirmPassword}
                onChange={(e) => handleInputChange('confirmPassword', e.target.value)}
                className={`mt-1 transition-all duration-300 focus:scale-105 hover:shadow-md ${errors.confirmPassword ? 'border-red-500 animate-pulse' : 'focus:border-purple-500'}`}
                placeholder="Confirm your password"
              />
              {errors.confirmPassword && (
                <p className="mt-1 text-sm text-red-600 animate-fade-in">{errors.confirmPassword}</p>
              )}
            </div>
          </div>

          <div className="animate-fade-in" style={{ animationDelay: '500ms' }}>
            <Button
              type="submit"
              className="w-full transform transition-all duration-300 hover:scale-105 hover:shadow-lg active:scale-95 bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-700 hover:to-pink-700 border-0 shadow-lg hover:shadow-xl text-white font-semibold py-3 relative overflow-hidden group"
              disabled={isLoading}
            >
              <div className="absolute inset-0 bg-gradient-to-r from-purple-500 to-pink-500 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
              <span className="relative z-10">
                {isLoading ? (
                  <span className="flex items-center justify-center space-x-2">
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    <span>Creating account...</span>
                  </span>
                ) : (
                  <span className="flex items-center justify-center space-x-2">
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z" />
                    </svg>
                    <span>Create Account</span>
                  </span>
                )}
              </span>
            </Button>
          </div>

          <div className="text-center animate-fade-in" style={{ animationDelay: '600ms' }}>
            <p className="text-sm text-gray-600">
              Already have an account?{' '}
              <Link to="/login" className="text-purple-600 hover:underline transition-all duration-300 hover:text-purple-800 relative after:content-[''] after:absolute after:w-full after:scale-x-0 after:h-0.5 after:bottom-0 after:left-0 after:bg-purple-600 after:origin-bottom-right after:transition-transform after:duration-300 hover:after:scale-x-100 hover:after:origin-bottom-left">
                Sign in here
              </Link>
            </p>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Signup;
